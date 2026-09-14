package server

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"ntocean-feast/backend/internal/config"
)

const sessionCookie = "kg_admin_session"

type Server struct {
	db       *sql.DB
	cfg      config.Config
	sessions *sessionStore
}

type sessionStore struct {
	mu     sync.Mutex
	tokens map[string]time.Time
}

type filters struct {
	Query    string
	Name     string
	Phone    string
	Email    string
	Code     string
	DateFrom string
	DateTo   string
	Status   string
	Page     int
	PageSize int
}

type codeRow struct {
	ID             uint64     `json:"id"`
	Code           string     `json:"code"`
	Registered     bool       `json:"registered"`
	Name           string     `json:"name"`
	Phone          string     `json:"phone"`
	Email          string     `json:"email"`
	TicketSource   string     `json:"ticket_source"`
	RestaurantName string     `json:"restaurant_name"`
	Satisfaction   *int       `json:"satisfaction"`
	Suggestion     string     `json:"suggestion"`
	RegisteredAt   *time.Time `json:"registered_at"`
}

type restaurantRow struct {
	ID       uint64 `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func New(db *sql.DB, cfg config.Config) http.Handler {
	s := &Server{db: db, cfg: cfg, sessions: &sessionStore{tokens: make(map[string]time.Time)}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/restaurants", s.listRestaurants)
	mux.HandleFunc("POST /api/registrations", s.register)
	mux.HandleFunc("POST /api/admin/login", s.login)
	mux.Handle("POST /api/admin/logout", s.requireAdmin(http.HandlerFunc(s.logout)))
	mux.Handle("GET /api/admin/session", s.requireAdmin(http.HandlerFunc(s.adminSession)))
	mux.Handle("GET /api/admin/summary", s.requireAdmin(http.HandlerFunc(s.summary)))
	mux.Handle("GET /api/admin/satisfaction", s.requireAdmin(http.HandlerFunc(s.satisfactionSummary)))
	mux.Handle("GET /api/admin/codes", s.requireAdmin(http.HandlerFunc(s.listCodes)))
	mux.Handle("GET /api/admin/export.csv", s.requireAdmin(http.HandlerFunc(s.exportCSV)))
	return recoverMiddleware(logMiddleware(mux))
}

func (s *Server) listRestaurants(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.QueryContext(r.Context(), `
		SELECT id, category, name, location
		FROM activity_restaurants WHERE active = TRUE
		ORDER BY sort_order ASC, id ASC`)
	if err != nil {
		log.Printf("list restaurants: %v", err)
		writeError(w, http.StatusInternalServerError, "無法讀取活動餐廳")
		return
	}
	defer rows.Close()
	items := make([]restaurantRow, 0)
	for rows.Next() {
		var item restaurantRow
		if err := rows.Scan(&item.ID, &item.Category, &item.Name, &item.Location); err != nil {
			writeError(w, http.StatusInternalServerError, "無法讀取活動餐廳")
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "無法讀取活動餐廳")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "資料庫尚未就緒")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string  `json:"name"`
		Phone        string  `json:"phone"`
		Email        string  `json:"email"`
		Code         string  `json:"code"`
		TicketSource string  `json:"ticket_source"`
		RestaurantID *uint64 `json:"restaurant_id"`
		Satisfaction int     `json:"satisfaction"`
		Suggestion   string  `json:"suggestion"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "請提供正確的登記資料")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	input.TicketSource = strings.TrimSpace(input.TicketSource)
	input.Suggestion = strings.TrimSpace(input.Suggestion)
	if message := validateRegistration(input.Name, input.Phone, input.Email, input.Code, input.TicketSource, input.RestaurantID, input.Satisfaction, input.Suggestion); message != "" {
		writeError(w, http.StatusUnprocessableEntity, message)
		return
	}
	if input.TicketSource == "partner_restaurant" {
		var active bool
		if err := s.db.QueryRowContext(r.Context(), `SELECT active FROM activity_restaurants WHERE id = ?`, *input.RestaurantID).Scan(&active); err != nil || !active {
			writeError(w, http.StatusUnprocessableEntity, "請選擇有效的活動合作餐廳")
			return
		}
	} else {
		input.RestaurantID = nil
	}

	result, err := s.db.ExecContext(r.Context(), `
		INSERT INTO registrations (
			lottery_code_id, name, phone, email, ticket_source,
			activity_restaurant_id, satisfaction, suggestion
		)
		SELECT id, ?, ?, ?, ?, ?, ?, ? FROM lottery_codes WHERE code = ?`,
		input.Name, input.Phone, input.Email, input.TicketSource,
		input.RestaurantID, input.Satisfaction, input.Suggestion, input.Code)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			writeError(w, http.StatusConflict, "此登錄編號已完成登記，無法重複使用")
			return
		}
		log.Printf("create registration: %v", err)
		writeError(w, http.StatusInternalServerError, "系統忙碌中，請稍後再試")
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "找不到此登錄編號，請確認後重新輸入")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "登記成功，祝您中獎！",
		"code":    input.Code,
	})
}

func validateRegistration(name, phone, email, code, ticketSource string, restaurantID *uint64, satisfaction int, suggestion string) string {
	if name == "" || phone == "" || email == "" || code == "" {
		return "姓名、聯絡電話、信箱與登錄編號皆為必填"
	}
	if ticketSource != "fuji_banquet" && ticketSource != "guihou_fair" && ticketSource != "partner_restaurant" {
		return "請選擇抽獎券來源"
	}
	if ticketSource == "partner_restaurant" && (restaurantID == nil || *restaurantID == 0) {
		return "請選擇活動合作餐廳"
	}
	if satisfaction < 1 || satisfaction > 5 {
		return "請選擇活動整體滿意度"
	}
	if suggestion == "" {
		return "請留下您對活動的寶貴建議"
	}
	if len([]rune(name)) > 100 || len(phone) > 30 || len(email) > 254 || len(code) > 32 || len([]rune(suggestion)) > 1000 {
		return "欄位內容過長，請確認後重新輸入"
	}
	address, err := mail.ParseAddress(email)
	if err != nil || !strings.EqualFold(address.Address, email) {
		return "請輸入有效的電子信箱"
	}
	return ""
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "請輸入帳號與密碼")
		return
	}
	userOK := subtle.ConstantTimeCompare([]byte(input.Username), []byte(s.cfg.AdminUser)) == 1
	passwordOK := subtle.ConstantTimeCompare([]byte(input.Password), []byte(s.cfg.AdminPassword)) == 1
	if !userOK || !passwordOK {
		writeError(w, http.StatusUnauthorized, "帳號或密碼錯誤")
		return
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		writeError(w, http.StatusInternalServerError, "無法建立登入狀態")
		return
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	expires := time.Now().Add(s.cfg.SessionLifetime)
	s.sessions.add(token, expires)
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: token, Path: "/", HttpOnly: true,
		Secure: s.cfg.CookieSecure, SameSite: http.SameSiteStrictMode, Expires: expires,
	})
	writeJSON(w, http.StatusOK, map[string]string{"username": input.Username})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		s.sessions.remove(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", HttpOnly: true, MaxAge: -1, SameSite: http.SameSiteStrictMode})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) adminSession(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"username": s.cfg.AdminUser})
}

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	var total, registered int
	err := s.db.QueryRowContext(r.Context(), `
		SELECT COUNT(*), COUNT(r.id)
		FROM lottery_codes c LEFT JOIN registrations r ON r.lottery_code_id = c.id`).Scan(&total, &registered)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "無法讀取統計資料")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"total": total, "registered": registered, "available": total - registered})
}

func (s *Server) satisfactionSummary(w http.ResponseWriter, r *http.Request) {
	counts := [5]int{}
	rows, err := s.db.QueryContext(r.Context(), `
		SELECT satisfaction, COUNT(*)
		FROM registrations
		WHERE satisfaction BETWEEN 1 AND 5
		GROUP BY satisfaction`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "無法讀取滿意度統計")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rating, count int
		if err := rows.Scan(&rating, &count); err != nil {
			writeError(w, http.StatusInternalServerError, "無法讀取滿意度統計")
			return
		}
		counts[rating-1] = count
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "無法讀取滿意度統計")
		return
	}
	total, weighted := 0, 0
	for index, count := range counts {
		total += count
		weighted += (index + 1) * count
	}
	average := 0.0
	if total > 0 {
		average = float64(weighted) / float64(total)
	}
	ratings := make([]map[string]any, 0, 5)
	for rating, count := range counts {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) * 100 / float64(total)
		}
		ratings = append(ratings, map[string]any{
			"rating": rating + 1, "count": count, "percentage": percentage,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"total": total, "average": average, "ratings": ratings})
}

func (s *Server) listCodes(w http.ResponseWriter, r *http.Request) {
	f, err := parseFilters(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	where, args := buildWhere(f)
	var total int
	countQuery := `SELECT COUNT(*) FROM lottery_codes c LEFT JOIN registrations r ON r.lottery_code_id = c.id` + where
	if err := s.db.QueryRowContext(r.Context(), countQuery, args...).Scan(&total); err != nil {
		log.Printf("count codes: %v", err)
		writeError(w, http.StatusInternalServerError, "無法讀取序號資料")
		return
	}

	query := `SELECT c.id, c.code, r.name, r.phone, r.email, r.ticket_source,
		ar.name, r.satisfaction, r.suggestion, r.registered_at
		FROM lottery_codes c
		LEFT JOIN registrations r ON r.lottery_code_id = c.id
		LEFT JOIN activity_restaurants ar ON ar.id = r.activity_restaurant_id` + where +
		` ORDER BY c.id ASC LIMIT ? OFFSET ?`
	queryArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	rows, err := s.db.QueryContext(r.Context(), query, queryArgs...)
	if err != nil {
		log.Printf("list codes: %v", err)
		writeError(w, http.StatusInternalServerError, "無法讀取序號資料")
		return
	}
	defer rows.Close()
	items, err := scanCodeRows(rows)
	if err != nil {
		log.Printf("scan codes: %v", err)
		writeError(w, http.StatusInternalServerError, "無法讀取序號資料")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items, "total": total, "page": f.Page, "page_size": f.PageSize,
		"total_pages": (total + f.PageSize - 1) / f.PageSize,
	})
}

func (s *Server) exportCSV(w http.ResponseWriter, r *http.Request) {
	f, err := parseFilters(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	where, args := buildWhere(f)
	query := `SELECT c.id, c.code, r.name, r.phone, r.email, r.ticket_source,
		ar.name, r.satisfaction, r.suggestion, r.registered_at
		FROM lottery_codes c
		LEFT JOIN registrations r ON r.lottery_code_id = c.id
		LEFT JOIN activity_restaurants ar ON ar.id = r.activity_restaurant_id` + where + ` ORDER BY c.id ASC`
	rows, err := s.db.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "無法匯出報表")
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="registrations-%s.csv"`, time.Now().Format("20060102-150405")))
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
	csvWriter := csv.NewWriter(w)
	_ = csvWriter.Write([]string{"序號", "登記狀態", "姓名", "聯絡電話", "信箱", "抽獎券來源", "活動合作餐廳", "滿意度", "活動建議", "登記日期"})
	items, err := scanCodeRows(rows)
	if err != nil {
		log.Printf("export scan: %v", err)
		return
	}
	for _, item := range items {
		status, registeredAt := "未登記", ""
		if item.Registered {
			status = "已登記"
			registeredAt = item.RegisteredAt.Format("2006-01-02 15:04:05")
		}
		satisfaction := ""
		if item.Satisfaction != nil {
			satisfaction = strconv.Itoa(*item.Satisfaction)
		}
		_ = csvWriter.Write([]string{
			safeCSV(item.Code), status, safeCSV(item.Name), safeCSV(item.Phone), safeCSV(item.Email),
			ticketSourceLabel(item.TicketSource), safeCSV(item.RestaurantName), satisfaction,
			safeCSV(item.Suggestion), registeredAt,
		})
	}
	csvWriter.Flush()
}

func safeCSV(value string) string {
	if value != "" && strings.ContainsRune("=+-@", rune(value[0])) {
		return "'" + value
	}
	return value
}

func ticketSourceLabel(value string) string {
	switch value {
	case "fuji_banquet":
		return "富基海派宴"
	case "guihou_fair":
		return "龜吼園遊會"
	case "partner_restaurant":
		return "活動合作餐廳"
	default:
		return ""
	}
}

func parseFilters(r *http.Request) (filters, error) {
	q := r.URL.Query()
	f := filters{
		Query: strings.TrimSpace(q.Get("q")), Name: strings.TrimSpace(q.Get("name")),
		Phone: strings.TrimSpace(q.Get("phone")), Email: strings.TrimSpace(q.Get("email")),
		Code: strings.TrimSpace(q.Get("code")), DateFrom: strings.TrimSpace(q.Get("date_from")),
		DateTo: strings.TrimSpace(q.Get("date_to")), Status: strings.TrimSpace(q.Get("status")),
		Page: 1, PageSize: 20,
	}
	if value, err := strconv.Atoi(q.Get("page")); err == nil && value > 0 {
		f.Page = value
	}
	if value, err := strconv.Atoi(q.Get("page_size")); err == nil && value > 0 && value <= 100 {
		f.PageSize = value
	}
	for _, value := range []string{f.DateFrom, f.DateTo} {
		if value != "" {
			if _, err := time.Parse("2006-01-02", value); err != nil {
				return f, errors.New("日期格式必須為 YYYY-MM-DD")
			}
		}
	}
	if f.Status != "" && f.Status != "registered" && f.Status != "available" {
		return f, errors.New("登記狀態不正確")
	}
	return f, nil
}

func buildWhere(f filters) (string, []any) {
	conditions := make([]string, 0, 9)
	args := make([]any, 0, 12)
	like := func(column, value string) {
		if value != "" {
			conditions = append(conditions, column+" LIKE ?")
			args = append(args, "%"+value+"%")
		}
	}
	if f.Query != "" {
		conditions = append(conditions, `(c.code LIKE ? OR r.name LIKE ? OR r.phone LIKE ? OR r.email LIKE ?)`)
		value := "%" + f.Query + "%"
		args = append(args, value, value, value, value)
	}
	like("r.name", f.Name)
	like("r.phone", f.Phone)
	like("r.email", f.Email)
	like("c.code", f.Code)
	if f.DateFrom != "" {
		conditions = append(conditions, "r.registered_at >= ?")
		args = append(args, f.DateFrom)
	}
	if f.DateTo != "" {
		conditions = append(conditions, "r.registered_at < DATE_ADD(?, INTERVAL 1 DAY)")
		args = append(args, f.DateTo)
	}
	if f.Status == "registered" {
		conditions = append(conditions, "r.id IS NOT NULL")
	} else if f.Status == "available" {
		conditions = append(conditions, "r.id IS NULL")
	}
	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func scanCodeRows(rows *sql.Rows) ([]codeRow, error) {
	items := make([]codeRow, 0)
	for rows.Next() {
		var item codeRow
		var name, phone, email, ticketSource, restaurantName, suggestion sql.NullString
		var satisfaction sql.NullInt64
		var registeredAt sql.NullTime
		if err := rows.Scan(
			&item.ID, &item.Code, &name, &phone, &email, &ticketSource,
			&restaurantName, &satisfaction, &suggestion, &registeredAt,
		); err != nil {
			return nil, err
		}
		item.Registered = registeredAt.Valid
		item.Name, item.Phone, item.Email = name.String, phone.String, email.String
		item.TicketSource = ticketSource.String
		item.RestaurantName = restaurantName.String
		item.Suggestion = suggestion.String
		if satisfaction.Valid {
			value := int(satisfaction.Int64)
			item.Satisfaction = &value
		}
		if registeredAt.Valid {
			value := registeredAt.Time
			item.RegisteredAt = &value
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil || !s.sessions.valid(cookie.Value) {
			writeError(w, http.StatusUnauthorized, "請先登入後台")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *sessionStore) add(token string, expires time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, expiry := range s.tokens {
		if expiry.Before(now) {
			delete(s.tokens, key)
		}
	}
	s.tokens[token] = expires
}

func (s *sessionStore) valid(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expires, ok := s.tokens[token]
	if !ok || expires.Before(time.Now()) {
		delete(s.tokens, token)
		return false
	}
	return true
}

func (s *sessionStore) remove(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started).Round(time.Millisecond))
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic: %v", recovered)
				writeError(w, http.StatusInternalServerError, "系統發生錯誤")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ShutdownContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

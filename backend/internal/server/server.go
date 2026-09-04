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
	ID           uint64     `json:"id"`
	Code         string     `json:"code"`
	Registered   bool       `json:"registered"`
	Name         string     `json:"name"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	RegisteredAt *time.Time `json:"registered_at"`
}

func New(db *sql.DB, cfg config.Config) http.Handler {
	s := &Server{db: db, cfg: cfg, sessions: &sessionStore{tokens: make(map[string]time.Time)}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("POST /api/registrations", s.register)
	mux.HandleFunc("POST /api/admin/login", s.login)
	mux.Handle("POST /api/admin/logout", s.requireAdmin(http.HandlerFunc(s.logout)))
	mux.Handle("GET /api/admin/session", s.requireAdmin(http.HandlerFunc(s.adminSession)))
	mux.Handle("GET /api/admin/summary", s.requireAdmin(http.HandlerFunc(s.summary)))
	mux.Handle("GET /api/admin/codes", s.requireAdmin(http.HandlerFunc(s.listCodes)))
	mux.Handle("GET /api/admin/export.csv", s.requireAdmin(http.HandlerFunc(s.exportCSV)))
	return recoverMiddleware(logMiddleware(mux))
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
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "請提供正確的登記資料")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	if message := validateRegistration(input.Name, input.Phone, input.Email, input.Code); message != "" {
		writeError(w, http.StatusUnprocessableEntity, message)
		return
	}

	result, err := s.db.ExecContext(r.Context(), `
		INSERT INTO registrations (lottery_code_id, name, phone, email)
		SELECT id, ?, ?, ? FROM lottery_codes WHERE code = ?`,
		input.Name, input.Phone, input.Email, input.Code)
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

func validateRegistration(name, phone, email, code string) string {
	if name == "" || phone == "" || email == "" || code == "" {
		return "姓名、聯絡電話、信箱與登錄編號皆為必填"
	}
	if len([]rune(name)) > 100 || len(phone) > 30 || len(email) > 254 || len(code) > 32 {
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

	query := `SELECT c.id, c.code, r.name, r.phone, r.email, r.registered_at
		FROM lottery_codes c LEFT JOIN registrations r ON r.lottery_code_id = c.id` + where +
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
	query := `SELECT c.id, c.code, r.name, r.phone, r.email, r.registered_at
		FROM lottery_codes c LEFT JOIN registrations r ON r.lottery_code_id = c.id` + where + ` ORDER BY c.id ASC`
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
	_ = csvWriter.Write([]string{"序號", "登記狀態", "姓名", "聯絡電話", "信箱", "登記日期"})
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
		_ = csvWriter.Write([]string{safeCSV(item.Code), status, safeCSV(item.Name), safeCSV(item.Phone), safeCSV(item.Email), registeredAt})
	}
	csvWriter.Flush()
}

func safeCSV(value string) string {
	if value != "" && strings.ContainsRune("=+-@", rune(value[0])) {
		return "'" + value
	}
	return value
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
		var name, phone, email sql.NullString
		var registeredAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Code, &name, &phone, &email, &registeredAt); err != nil {
			return nil, err
		}
		item.Registered = registeredAt.Valid
		item.Name, item.Phone, item.Email = name.String, phone.String, email.String
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

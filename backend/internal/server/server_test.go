package server

import (
	"strings"
	"testing"
	"time"
)

func TestValidateRegistration(t *testing.T) {
	restaurantID := uint64(1)
	tests := []struct {
		name, phone, email, code, source, suggestion string
		restaurantID                                 *uint64
		satisfaction                                 int
		wantError                                    bool
	}{
		{"王小明", "0912345678", "user@example.com", "ABCDE00001", "fuji_banquet", "很棒", nil, 5, false},
		{"王小明", "0912345678", "user@example.com", "ABCDE00001", "partner_restaurant", "很棒", &restaurantID, 4, false},
		{"", "0912345678", "user@example.com", "ABCDE00001", "fuji_banquet", "很棒", nil, 5, true},
		{"王小明", "0912345678", "not-an-email", "ABCDE00001", "fuji_banquet", "很棒", nil, 5, true},
		{"王小明", "0912345678", "user@example.com", "ABCDE00001", "partner_restaurant", "很棒", nil, 5, true},
		{"王小明", "0912345678", "user@example.com", "ABCDE00001", "guihou_fair", "", nil, 6, true},
	}
	for _, test := range tests {
		got := validateRegistration(
			test.name, test.phone, test.email, test.code, test.source,
			test.restaurantID, test.satisfaction, test.suggestion,
		)
		if (got != "") != test.wantError {
			t.Fatalf("validateRegistration() = %q, wantError %v", got, test.wantError)
		}
	}
}

func TestBuildWhere(t *testing.T) {
	where, args := buildWhere(filters{Query: "王", Status: "registered", DateFrom: "2026-09-01"})
	for _, fragment := range []string{"c.code LIKE", "r.id IS NOT NULL", "r.registered_at >="} {
		if !strings.Contains(where, fragment) {
			t.Errorf("where %q does not contain %q", where, fragment)
		}
	}
	if len(args) != 5 {
		t.Errorf("len(args) = %d, want 5", len(args))
	}
}

func TestSessionExpiry(t *testing.T) {
	store := &sessionStore{tokens: make(map[string]time.Time)}
	store.add("valid", time.Now().Add(time.Minute))
	store.add("expired", time.Now().Add(-time.Minute))
	if !store.valid("valid") || store.valid("expired") {
		t.Fatal("session validity did not respect expiry")
	}
}

func TestSafeCSV(t *testing.T) {
	if got := safeCSV("=SUM(1,1)"); got != "'=SUM(1,1)" {
		t.Fatalf("safeCSV() = %q", got)
	}
	if got := safeCSV("王小明"); got != "王小明" {
		t.Fatalf("safeCSV() changed regular text to %q", got)
	}
}

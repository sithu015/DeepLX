package translate

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestTranslateCharacterLimitCheck(t *testing.T) {
	// Create a text that exceeds 1500 characters
	longText := strings.Repeat("a", 1501)

	// Case 1: Free request (dlSession == "") should be rejected immediately
	resFree, err := TranslateByDeepLX(context.Background(), "EN", "DE", longText, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resFree.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected Code 413 for free request exceeding limit, got %d", resFree.Code)
	}

	// Case 2: Pro request (dlSession != "") should bypass character limit check.
	// Since we pass a dummy session, it will try to hit the upstream oneshot-pro URL
	// and fail/timeout or return 401/403/503. The point is, it should NOT return 413.
	resPro, _ := TranslateByDeepLX(context.Background(), "EN", "DE", longText, "", "", "dummy_pro_session")
	if resPro.Code == http.StatusRequestEntityTooLarge {
		t.Errorf("expected Pro request to bypass 413 character limit check, but got 413")
	}
}

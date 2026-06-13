package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &Config{
		Token: "test-secret-token",
	}

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(authMiddleware(cfg))
		r.POST("/translate", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		r.POST("/:token/translate", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "path_token": c.Param("token")})
		})
		return r
	}

	r := setupRouter()

	tests := []struct {
		name           string
		method         string
		url            string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "Valid Bearer Header",
			method:         "POST",
			url:            "/translate",
			headers:        map[string]string{"Authorization": "Bearer test-secret-token"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Custom DeepL-Auth-Key Header",
			method:         "POST",
			url:            "/translate",
			headers:        map[string]string{"Authorization": "DeepL-Auth-Key test-secret-token"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Header Token",
			method:         "POST",
			url:            "/translate",
			headers:        map[string]string{"Authorization": "Bearer wrong-token"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Approach A - Valid Query Parameter token",
			method:         "POST",
			url:            "/translate?token=test-secret-token",
			headers:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Approach A - Valid Query Parameter key",
			method:         "POST",
			url:            "/translate?key=test-secret-token",
			headers:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Approach A - Invalid Query Parameter token",
			method:         "POST",
			url:            "/translate?token=wrong-token",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Approach B - Valid Path Parameter token",
			method:         "POST",
			url:            "/test-secret-token/translate",
			headers:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Approach B - Invalid Path Parameter token",
			method:         "POST",
			url:            "/wrong-token/translate",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Authentication",
			method:         "POST",
			url:            "/translate",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte("{}")))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

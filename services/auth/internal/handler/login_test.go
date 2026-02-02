package handler

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHashToken_Deterministic(t *testing.T) {
	token := "my-test-token-123"
	hash1 := hashToken(token)
	hash2 := hashToken(token)
	if hash1 != hash2 {
		t.Errorf("hashToken is not deterministic: %s != %s", hash1, hash2)
	}
}

func TestHashToken_DifferentInputs(t *testing.T) {
	hash1 := hashToken("token-a")
	hash2 := hashToken("token-b")
	if hash1 == hash2 {
		t.Error("hashToken produced same hash for different inputs")
	}
}

func TestHashToken_Length(t *testing.T) {
	hash := hashToken("any-token")
	// SHA256 produces 32 bytes = 64 hex characters
	if len(hash) != 64 {
		t.Errorf("expected hash length 64, got %d", len(hash))
	}
}

func TestIsRequestHTTPS_XForwardedProto(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("X-Forwarded-Proto", "https")

	if !isRequestHTTPS(c) {
		t.Error("expected isRequestHTTPS to return true with X-Forwarded-Proto: https")
	}
}

func TestIsRequestHTTPS_NoHTTPS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if isRequestHTTPS(c) {
		t.Error("expected isRequestHTTPS to return false without HTTPS indicators")
	}
}

func TestIsRequestHTTPS_DirectTLS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.TLS = &tls.ConnectionState{}

	if !isRequestHTTPS(c) {
		t.Error("expected isRequestHTTPS to return true with direct TLS")
	}
}

package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func newTestJWT() *jwt.JWT {
	return jwt.New([]byte("test-secret-key-for-testing-only"))
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRoleMiddleware_Allowed(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("role", uint(RoleAdmin))
		c.Next()
	})
	r.Use(RoleMiddleware(RoleAdmin))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRoleMiddleware_Denied(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("role", uint(RoleUser))
		c.Next()
	})
	r.Use(RoleMiddleware(RoleAdmin))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRoleMiddleware_NoRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// No role set in context
	r.Use(RoleMiddleware(RoleAdmin))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRoleMiddleware_MultipleRoles(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("role", uint(RoleManager))
		c.Next()
	})
	r.Use(RoleMiddleware(RoleAdmin, RoleManager))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_Valid(t *testing.T) {
	j := newTestJWT()

	tokenData, err := j.GenerateAccessToken("test@example.com", RoleAdmin, "access")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	var gotEmail string
	var gotRole uint

	r.Use(AuthMiddleware(j))
	r.GET("/test", func(c *gin.Context) {
		email, _ := c.Get("email")
		gotEmail = email.(string)
		role, _ := c.Get("role")
		gotRole = role.(uint)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.Token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if gotEmail != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", gotEmail)
	}
	if gotRole != RoleAdmin {
		t.Errorf("expected role %d, got %d", RoleAdmin, gotRole)
	}
}

func TestAuthMiddleware_Missing(t *testing.T) {
	j := newTestJWT()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(AuthMiddleware(j))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	j := newTestJWT()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(AuthMiddleware(j))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_RefreshAsAccess(t *testing.T) {
	j := newTestJWT()

	// Generate a refresh token
	tokenData, err := j.GenerateRefreshToken("test@example.com", RoleUser)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(AuthMiddleware(j))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.Token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	// Should indicate invalid token
	if resp["message"] == nil {
		t.Error("expected error message in response")
	}
}

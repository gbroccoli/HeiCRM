package handler

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLogActivity_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO user_activity_log").
		WithArgs(uint64(1), "test_action", []byte(`{"key":"value"}`), "192.0.2.1", "TestAgent/1.0").
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	c.Request.Header.Set("User-Agent", "TestAgent/1.0")
	c.Request.RemoteAddr = "192.0.2.1:1234"

	logActivity(db, 1, "test_action", map[string]string{"key": "value"}, c)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestLogActivity_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO user_activity_log").
		WillReturnError(sqlmock.ErrCancelled)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)

	// Should not panic
	logActivity(db, 1, "test_action", map[string]string{"key": "value"}, c)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestLogActivity_MarshalError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)

	// math.Inf(1) cannot be marshaled to JSON
	// Should not panic — returns early after marshal error
	logActivity(db, 1, "test_action", math.Inf(1), c)
}

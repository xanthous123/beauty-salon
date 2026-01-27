package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// tests

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - within limit", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		limit := 2
		window := time.Minute

		r := gin.New()
		r.Use(RateLimiter(db, limit, window))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		mock.ExpectIncr("rate_limit:127.0.0.1").SetVal(1)
		mock.ExpectExpire("rate_limit:127.0.0.1", window).SetVal(true)

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure - limit exceeded", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		limit := 1

		r := gin.New()
		r.Use(RateLimiter(db, limit, time.Minute))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		mock.ExpectIncr("rate_limit:127.0.0.1").SetVal(2)

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusTooManyRequests, resp.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail Open - redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		r := gin.New()
		r.Use(RateLimiter(db, 5, time.Minute))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		mock.ExpectIncr("rate_limit:127.0.0.1").SetErr(redis.ErrClosed)

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

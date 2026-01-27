package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// tests
func TestAuthMiddleware(t *testing.T) {
	secret := "test_secret"
	os.Setenv("JWT_SECRET", secret)
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, r := gin.CreateTestContext(resp)

		r.Use(AuthMiddleware())
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(resp, c.Request)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Authorization required")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, r := gin.CreateTestContext(resp)

		r.Use(AuthMiddleware())
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(resp, c.Request)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid token")
	})

	t.Run("Valid Token Success", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, r := gin.CreateTestContext(resp)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": float64(123),
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		r.Use(AuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, uint(123), userID)
			c.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)
		r.ServeHTTP(resp, c.Request)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

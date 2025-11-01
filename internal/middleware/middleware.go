package middleware

import (
	"ToDo/internal/model"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	db *gorm.DB
}

func NewAuthMiddleware(database *gorm.DB) AuthMiddleware {
	return AuthMiddleware{db: database}
}

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing auth cookie")
		}

		tokenString := cookie.Value

		if !strings.HasPrefix(tokenString, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing bearer token")
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unxpected signing method: %v", token.Header)
			}

			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: "+err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			// Check exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token expired")
			}

			// Find user with token sub
			var user model.User
			err := m.db.First(&user, "login = ?", claims["sub"]).Error
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
			}

			c.Set("user", user)

		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		// Continue
		return next(c)
	}
}

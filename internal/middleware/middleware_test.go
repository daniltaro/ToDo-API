package middleware

import (
	"ToDo/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMiddleware_auth(t *testing.T) {
	// Create bd
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	db.Create(&model.User{Login: "testUser"})

	// Create middleware
	authMiddleware := NewAuthMiddleware(db)
	os.Setenv("SECRET", "test_secret")

	// Handler after middleware
	nextHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}

	e := echo.New()

	testTable := []struct {
		name               string
		expectedStatusCode int
		tokenGen           func() string
	}{
		{
			name:               "OK",
			expectedStatusCode: 200,
			tokenGen: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "testUser",
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte("test_secret"))
				return "Bearer " + tokenStr
			},
		},
		{
			name:               "Expired token",
			expectedStatusCode: 401,
			tokenGen: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "testUser",
					"exp": time.Now().Add(-time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte("test_secret"))
				return "Bearer " + tokenStr
			},
		},
		{
			name:               "Invalid login",
			expectedStatusCode: 401,
			tokenGen: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "iminvalid",
					"exp": time.Now().Add(-time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte("test_secret"))
				return "Bearer " + tokenStr
			},
		},
		{
			name:               "No cookie",
			expectedStatusCode: 401,
			tokenGen:           func() string { return "" },
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			c := e.NewContext(req, w)

			if testCase.tokenGen() != "" {
				cookie := new(http.Cookie)
				cookie.Name = "Authorization"
				cookie.Value = testCase.tokenGen()
				req.AddCookie(cookie)
			}

			handler := authMiddleware.RequireAuth(nextHandler)
			err := handler(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}

}

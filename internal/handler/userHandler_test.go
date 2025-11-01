package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"ToDo/internal/model"
	mock_service "ToDo/internal/service/mocks"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_signup(t *testing.T) {
	type mockBehavior func(s *mock_service.MockService, user model.User)

	test_table := []struct {
		name                string
		inputBody           string
		inputUser           model.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"login":"test","password":"test"}`,
			inputUser: model.User{
				Login:    "test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockService, user model.User) {
				s.EXPECT().AddUser(&user).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"success\":\"User added\"}\n",
		},
		{
			name:                "Empty fields(wrong input)",
			inputBody:           "",
			inputUser:           model.User{},
			mockBehavior:        func(s *mock_service.MockService, user model.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"error\":\"Invalid request\"}\n",
		},
		{
			name:      "Service failure",
			inputBody: `{"login":"test","password":"test"}`,
			inputUser: model.User{
				Login:    "test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockService, user model.User) {
				s.EXPECT().AddUser(&user).Return(errors.New("some server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"error\":\"Could not add user\"}\n",
		},
	}

	for _, testCase := range test_table {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockService(c)
			testCase.mockBehavior(service, testCase.inputUser)

			handler := NewUserHandler(service)

			e := echo.New()
			e.POST("/signup", handler.Signup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestUserHandler_login(t *testing.T) {
	type mockBehavior func(s *mock_service.MockService, user model.User)
	os.Setenv("SECRET", "test_secret")

	test_table := []struct {
		name               string
		inputBody          string
		inputUser          model.User
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedCookie     *http.Cookie
	}{
		{
			name:      "OK",
			inputBody: `{"login":"test","password":"test"}`,
			inputUser: model.User{
				Login:    "test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockService, user model.User) {
				s.EXPECT().LookUpReqUser(&user).Return(nil)
			},
			expectedStatusCode: 200,
			expectedCookie: &http.Cookie{
				Name:  "Authorization",
				Value: "Bearer " + generateToken("test", time.Now().Add(time.Hour*24*30).Unix()),
			},
		},
		{
			name:      "Invalid user",
			inputBody: `{"login":"invalid","password":"invalid"}`,
			inputUser: model.User{
				Login:    "invalid",
				Password: "invalid",
			},
			mockBehavior: func(s *mock_service.MockService, user model.User) {
				s.EXPECT().LookUpReqUser(&user).Return(fmt.Errorf("wrong password or login"))
			},
			expectedStatusCode: 400,
			expectedCookie:     &http.Cookie{},
		},
	}
	for _, testCase := range test_table {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockService(c)
			testCase.mockBehavior(service, testCase.inputUser)

			handler := NewUserHandler(service)

			e := echo.New()
			e.POST("/login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			actualCookie := &http.Cookie{}
			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == "Authorization" {
					actualCookie = cookie
					break
				}
			}
			assert.Equal(t, testCase.expectedCookie.Value, actualCookie.Value)
		})
	}
}

func generateToken(sub string, exp int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString
}

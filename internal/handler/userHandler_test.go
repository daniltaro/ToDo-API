package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"ToDo/internal/model"
	mock_service "ToDo/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signup(t *testing.T) {
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

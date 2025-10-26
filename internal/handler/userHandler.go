package handler

import (
	"net/http"
	"os"
	"time"

	"ToDo/internal/model"
	"ToDo/internal/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type UserHandler struct {
	s service.Service
}

func NewUserHandler(service service.Service) UserHandler {
	return UserHandler{s: service}
}

func (h *UserHandler) Signup(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := h.s.AddUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not add user",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "User added"})
}

func (h *UserHandler) Login(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := h.s.LookUpReqUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Wrong password or login",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Login,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create token",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

func (h *UserHandler) Validate(c echo.Context) error {
	user := c.Get("user")
	return c.JSON(http.StatusOK, user.(model.User))
}

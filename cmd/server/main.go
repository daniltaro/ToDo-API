package main

import (
	"log"

	"ToDo/internal/handler"
	"ToDo/internal/initializers"
	"ToDo/internal/middleware"
	"ToDo/internal/repository"
	"ToDo/internal/service"

	"github.com/labstack/echo"
	echoMiddlware "github.com/labstack/echo/middleware"
)

func main() {
	err := initializers.InitENV()
	if err != nil {
		log.Fatalf("Could not find .env: %v", err)
	}
	database, err := initializers.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	repo := repository.NewRepository(database)
	authMiddleware := middleware.NewAuthMiddleware(database)
	service := service.NewService(repo)
	taskHandler := handler.NewTaskHandler(service)
	userHandler := handler.NewUserHandler(service)

	e := echo.New()

	e.GET("/tasks", taskHandler.GetTasks, authMiddleware.RequireAuth)
	e.POST("/tasks", taskHandler.PostTasks, authMiddleware.RequireAuth)
	e.PATCH("/tasks/:id", taskHandler.PatchTasks, authMiddleware.RequireAuth)
	e.DELETE("/tasks/:id", taskHandler.DeleteTasks, authMiddleware.RequireAuth)
	e.GET("/validate", userHandler.Validate, authMiddleware.RequireAuth)

	e.Use(echoMiddlware.Logger())

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)

	e.Start("localhost:8080")
}

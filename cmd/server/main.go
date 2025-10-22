package main

import (
	"log"

	"ToDo/internal/db"
	"ToDo/internal/repository"
	"ToDo/internal/service"
	"ToDo/internal/handler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)


func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	repo := repository.NewTaskRepository(database)
	service := service.NewTaskService(repo)
	handler := handler.NewTaskHandler(service)

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/tasks", handler.GetTasks)
	e.POST("/tasks", handler.PostTasks)
	e.PATCH("/tasks/:id", handler.PatchTasks)

	e.Start("localhost:8080")
}

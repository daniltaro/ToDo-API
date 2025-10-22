package handler

import (
	"net/http"

	"ToDo/internal/model"
	"ToDo/internal/service"

	"github.com/labstack/echo"
)

type TaskHandler struct {
	s service.TaskService
}

func NewTaskHandler(service service.TaskService) TaskHandler {
	return TaskHandler{s: service}
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	tasks, err := h.s.GetAllTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get tasks"})
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) PostTasks(c echo.Context) error {
	var task model.Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task"})
	}

	if err := h.s.CreateTask(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not create task"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Task added"})
}

func (h *TaskHandler) PatchTasks(c echo.Context) error {
	id := c.Param("id")
	var taskCondition model.TaskCondittion
	if err := c.Bind(&taskCondition); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := h.s.ChangeTaskCondition(id, taskCondition.IsDone); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not change task condition"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Task condition changed"})
}

func (h *TaskHandler) DeleteTasks(c echo.Context) error {
	id := c.Param("id")
	if err := h.s.DeleteTask(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete task"})
	}
	return c.NoContent(http.StatusNoContent)
}

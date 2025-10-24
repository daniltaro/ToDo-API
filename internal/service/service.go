package service

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"ToDo/internal/model"
	"ToDo/internal/repository"
)

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	CreateTask(task *model.Task) error
	ChangeTaskCondition(id string, isDone bool) error
	DeleteTask(id string) error
}

type TaskServ struct {
	r repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServ{r: repo}
}

func (s *TaskServ) GetAllTasks() ([]model.Task, error) {
	return s.r.GetAllTasks()
}

func (s *TaskServ) CreateTask(task *model.Task) error {
	task_ := task
	task_.ID = uuid.New()
	return s.r.CreateTask(task_)
}

func (s *TaskServ) ChangeTaskCondition(id string, isDone bool) error {
	task, err := s.r.GetTaskByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetTaskByID: %v", err)
		return err
	}

	task.IsDone = isDone

	if err = s.r.UpdateTask(&task); err != nil {
		fmt.Fprintf(os.Stderr, "UpdateTask: %v", err)
		return err
	}

	return nil
}

func (s *TaskServ) DeleteTask(id string) error {
	return s.r.DeleteTask(id)
}

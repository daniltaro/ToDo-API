package service

import (
	"ToDo/internal/model"
	"ToDo/internal/repository"

	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	CreateTask(task *model.Task) error
	ChangeTaskCondition(id string, isDone bool) error
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
	task_.ID = uuid.NewString()
	return s.r.CreateTask(task_)
}

func (s *TaskServ) ChangeTaskCondition(id string, isDone bool) error {
	task, err := s.r.GetTaskByID(id)
	if err != nil {
		return err
	}

	task.IsDone = isDone

	if err = s.r.UpdateTask(&task); err != nil {
		return err
	}

	return nil
}

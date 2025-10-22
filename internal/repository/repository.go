package repository

import (
	"ToDo/internal/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetAllTasks() ([]model.Task, error)
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	GetTaskByID(id string) (model.Task, error)
	DeleteTask(id string) error
}

type TaskRep struct {
	db *gorm.DB
}

func NewTaskRepository(database *gorm.DB) TaskRepository {
	return &TaskRep{db: database}
}

func (r *TaskRep) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TaskRep) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRep) UpdateTask(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRep) GetTaskByID(id string) (model.Task, error) {
	var task model.Task
	err := r.db.First(&task, "id = ?", id).Error
	return task, err
}

func (r *TaskRep) DeleteTask(id string) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}

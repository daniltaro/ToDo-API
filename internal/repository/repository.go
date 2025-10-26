package repository

import (
	"ToDo/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetAllTasks(login string) ([]model.Task, error)
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	GetTaskByID(id string) (model.Task, error)
	DeleteTask(id, login string) error
	AddUser(user *model.User) error
	GetUserByLogin(login string) (model.User, error)
}

type Rep struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &Rep{db: database}
}

func (r *Rep) GetAllTasks(login string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("login = ?", login).Find(&tasks).Error
	return tasks, err
}

func (r *Rep) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *Rep) UpdateTask(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *Rep) GetTaskByID(id string) (model.Task, error) {
	var task model.Task
	err := r.db.First(&task, "id = ?", id).Error
	return task, err
}

func (r *Rep) DeleteTask(id, login string) error {
	return r.db.Where("login = ?", login).Delete(&model.Task{}, "id = ?", id).Error
}

func (r *Rep) AddUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Rep) GetUserByLogin(login string) (model.User, error) {
	var user model.User
	err := r.db.First(&user, "login = ?", login).Error
	return user, err
}

package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/daniltaro/ToDo-API/internal/model"
	"github.com/daniltaro/ToDo-API/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service interface {
	GetAllTasks(login string) ([]model.Task, error)
	CreateTask(task *model.Task, login string) error
	ChangeTaskCondition(id, login string, isDone bool) error
	DeleteTask(id, login string) error
	AddUser(user *model.User) error
	LookUpReqUser(user *model.User) error
}

type Serv struct {
	r repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &Serv{r: repo}
}

func (s *Serv) GetAllTasks(login string) ([]model.Task, error) {
	return s.r.GetAllTasks(login)
}

func (s *Serv) CreateTask(task *model.Task, login string) error {
	task.ID = uuid.New()
	task.Login = login
	return s.r.CreateTask(task)
}

func (s *Serv) ChangeTaskCondition(id, login string, isDone bool) error {
	task, err := s.r.GetTaskByID(id)
	if err != nil {
		return err
	}
	if task.Login != login {
		return errors.New("could not find task")
	}

	task.IsDone = isDone

	if err = s.r.UpdateTask(&task); err != nil {
		return err
	}

	return nil
}

func (s *Serv) DeleteTask(id, login string) error {
	return s.r.DeleteTask(id, login)
}

func (s *Serv) AddUser(user *model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	return s.r.AddUser(&model.User{
		Password:  string(hash),
		Login:     user.Login,
		CreatedAt: time.Now(),
	})
}

// Checks password and login in db
func (s *Serv) LookUpReqUser(user *model.User) error {
	var body model.User
	var err error
	if body, err = s.r.GetUserByLogin(user.Login); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(body.Password), []byte(user.Password))
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.UserCreator) (int, string, error)
	GenerateToken(username, password string) (model.UserGet, string, error)
	ConfirmUser(confirm string) (int, error)
	GetConfirm(login string) (model.UserGet, error)
	ResetPassword(login, password string) error
	RefreshToken(reftoken string, preId int) (int, string, error)
	ParseToken(token string) (int, error)
}
type Users interface {
	MeUser(id int) (model.UserCreator, error)
	GetUserAll() ([]model.UserCreator, error)
	UpdateLink(arr []model.Link, id int) error
	UpdateUser(id int, newData model.UserCreator) (model.UserCreator, error)
}

type Service struct {
	Authorization
	Users
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewAuthService(repos.Authorization),
		NewUserService(repos.Users),
	}
}

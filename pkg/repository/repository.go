package repository

import (
	"github.com/Constantine-Ka/user-service/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.UserCreator) (int, string, error)
	GetUser(login, email, password string, id int) (model.UserGet, error)
	ConfirmUser(confirm string, token model.RefreshData) (int, error)
	GetConfirm(login, column string) (model.UserGet, error)
	ResetPassword(data model.UserCreator) error
	RefreshToken(reftoken string) (int, error)
}
type Users interface {
	GetUserOne(id int) (model.UserGet, error)
	GetUserAll() ([]model.UserGet, error)
	UpdateLink(arr []model.Link, id int) error
	UpdateUser(id int, newData model.UserCreator) (model.UserGet, error)
}

type Repository struct {
	Authorization
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Users:         NewUserPostgres(db),
	}
}

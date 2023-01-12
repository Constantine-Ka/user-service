package repository

import (
	"errors"
	"fmt"
	"github.com/Constantine-Ka/user-service/tools/logging"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateUser(user model.UserCreator) (int, string, error) {
	logger := logging.GetLogger()
	logger.Info("CreateUser")
	var id int
	var confirmation string
	query := fmt.Sprintf("INSERT INTO %s (login, email, password_hash,confirmation, first_name, second_name, last_name, image, gender, birthday, description, register_date) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id,confirmation", userTable)
	row := r.db.QueryRow(query, user.Login, user.EMail, user.Password, user.Confirmation, user.FirstName, user.SecondName, user.LastName, user.ImagePath, user.Gender, user.Birthday.TimeStamp, user.Description, time.Now().Unix())
	if err := row.Scan(&id, &confirmation); err != nil {
		return 0, "", err
	}
	return id, confirmation, nil
}
func (r *AuthPostgres) GetUser(login, email, password string, id int) (model.UserGet, error) {
	var (
		user  model.UserGet
		query string
		err   error
	)
	logger := logging.GetLogger()
	if (login == "") && (email == "") && (id == 0) {
		return model.UserGet{}, errors.New("login and E-mail is Empty")
	} else if id != 0 {
		query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", userTable)
		err = r.db.Get(&user, query, id)
	} else if login == "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", userTable)
		err = r.db.Get(&user, query, email, password)
	} else if email == "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password_hash=$2", userTable)
		err = r.db.Get(&user, query, login, password)
	} else {
		err = errors.New("Я где-то проебался")
	}
	if id != 0 {
		user.Token.String = "is Secret 😛😛😛"
	}
	if err == errors.New("\u003cnil\u003e") {
		return model.UserGet{}, errors.New("аккаунт не активирован, обратитесь к администратору ресурса")
	}
	if !user.IsConfirmation.Bool {
		return user, errors.New("аккаунт не активирован. Проверьте почту,указанную при регистрации")
	}
	logger.Info(user.IsConfirmation)
	return user, err
}
func (r *AuthPostgres) ConfirmUser(confirm string, token model.RefreshData) (int, error) {
	var id int
	query := fmt.Sprintf("UPDATE %s SET is_confirm=true, refresh_token=$2, refresh_token_expires=$3 WHERE confirmation=$1 RETURNING ID", userTable)
	err := r.db.Get(&id, query, confirm, token.Token, token.ExpiresAt)
	logger := logging.GetLogger()
	logger.Error(err)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (r *AuthPostgres) GetConfirm(login, column string) (model.UserGet, error) {
	var response model.UserGet

	query := fmt.Sprintf("SELECT confirmation,email FROM %s WHERE %s=$1", userTable, column)
	logrus.Warn(query)
	var err = r.db.Get(&response, query, login)
	if err != nil {
		return response, err
	}
	return response, nil
}
func (r *AuthPostgres) ResetPassword(data model.UserCreator) error {
	var query string
	queryTemplate := "UPDATE %s SET password_hash=$2 WHERE %s=$1"
	logrus.Warn(data.Id)
	if data.Id != 0 {
		query = fmt.Sprintf("UPDATE %d SET password_hash=$2 WHERE %s=$1", userTable, "id")
		row := r.db.QueryRow(query, data.Id, data.Password)
		if err := row.Err(); err != nil {
			return err
		}
	} else if data.Login != "" {
		query = fmt.Sprintf(queryTemplate, userTable, "login")
		row := r.db.QueryRow(query, data.Login, data.Password)
		if err := row.Err(); err != nil {
			return err
		}
	} else if data.EMail != "" {
		query = fmt.Sprintf(queryTemplate, userTable, "email")
		row := r.db.QueryRow(query, data.EMail, data.Password)
		if err := row.Err(); err != nil {
			return err
		}
	} else if data.Confirmation != "" {
		query = fmt.Sprintf(queryTemplate, userTable, "confirmation")
		row := r.db.QueryRow(query, data.Confirmation, data.Password)
		if err := row.Err(); err != nil {
			return err
		}
	} else {
		return errors.New("id is empty")
	}
	return nil
}
func (r *AuthPostgres) UpdateUser(data model.UserCreator) error {
	var query string
	queryTemplate := "UPDATE %s SET login=$2, email=$3,first_name=$4,second_name=$5,last_name=$6,image=$7,gender=$8,birthday=$9,description=$10, my_links=$11 WHERE %s=$1"

	if data.Id != 0 {
		query = fmt.Sprintf(queryTemplate, userTable, "id")
		row := r.db.QueryRow(query, data.Id, data.Login, data.EMail, data.FirstName, data.SecondName, data.LastName, data.ImagePath, data.Gender, data.Birthday, data.Description, data.Links)
		if err := row.Err(); err != nil {
			return err
		}

	} else {
		return errors.New("id is empty")
	}

	return nil

}
func (r *AuthPostgres) RefreshToken(reftoken string) (int, error) {
	var user model.UserGet
	query := fmt.Sprintf("SELECT id, refresh_token_expires FROM %s WHERE refresh_token = $1", userTable)
	err := r.db.Get(&user, query, reftoken)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

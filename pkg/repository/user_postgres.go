package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/tools/logging"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// GetUserOne
func (r *UserPostgres) GetUserOne(id int) (model.UserGet, error) {
	var (
		user  model.UserGet
		query string
		err   error
	)
	logger := logging.GetLogger()
	if id != 0 {
		query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", userTable)
		err = r.db.Get(&user, query, id)
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

// GetUserAll
func (r *UserPostgres) GetUserAll() ([]model.UserGet, error) {
	var user []model.UserGet
	query := fmt.Sprintf("SELECT * FROM %s", userTable)
	err := r.db.Select(&user, query)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (r *UserPostgres) UpdateLink(arr []model.Link, id int) error {
	query := fmt.Sprintf("UPDATE %s SET my_links = $1 WHERE id = $2", userTable)
	bytes, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	rows, err := r.db.Query(query, bytes, id)
	if err != nil {
		return err
	}
	fmt.Println(rows)
	return nil
}

func (r UserPostgres) UpdateUser(id int, data model.UserCreator) (model.UserGet, error) {
	var user model.UserGet
	queryTemplate := fmt.Sprintf("UPDATE %s SET login=$2, email=$3,first_name=$4,second_name=$5,last_name=$6,image=$7,gender=$8,birthday=$9,description=$10, my_links=$11 WHERE %s=$1 RETURNING *", userTable, "id")

	bytes, err := json.Marshal(data.Links)
	if err != nil {
		return model.UserGet{}, err
	}
	err = r.db.Get(&user, queryTemplate, id, data.Login, data.EMail, data.FirstName, data.SecondName, data.LastName, data.ImagePath, data.Gender, data.Birthday.TimeStamp, data.Description, bytes)
	if err != nil {
		return model.UserGet{}, err
	}

	return user, nil
}

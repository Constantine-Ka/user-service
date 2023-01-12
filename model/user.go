package model

import (
	"database/sql"
)

type (
	Link struct {
		Title   string `json:"title,omitempty" binding:"required"`
		URL     string `json:"URL,omitempty" binding:"required"`
		Tooltip string `json:"tooltip,omitempty"`
	}
	UserGet struct {
		Id             int            `json:"id" db:"id"`
		Login          string         `json:"login" binding:"required"`
		EMail          string         `json:"email"`
		Password       string         `json:"password" binding:"required" db:"password_hash"`
		Confirmation   sql.NullString `json:"confirmation,omitempty" db:"confirmation"`
		IsConfirmation sql.NullBool   `json:"is_confirm,omitempty" db:"is_confirm"`
		Token          sql.NullString `json:"token" db:"refresh_token"`
		ExpiresAt      sql.NullInt64  `json:"exp_token" db:"refresh_token_expires"`
		FirstName      sql.NullString `json:"firstName,omitempty" db:"first_name"`
		SecondName     sql.NullString `json:"secondName,omitempty"db:"second_name"`
		LastName       sql.NullString `json:"lastName,omitempty" db:"last_name"`
		ImagePath      sql.NullString `json:"imagePath,omitempty" db:"image"`
		Gender         sql.NullInt16  `json:"gender,omitempty" db:"gender"`
		Birthday       sql.NullInt64  `json:"birthday,omitempty" db:"birthday"`
		Description    sql.NullString `json:"description,omitempty" db:"description"`
		RegisterDate   sql.NullInt64  `json:"registerDate,omitempty" db:"register_date"`
		LastDate       sql.NullInt64  `json:"lastDate,omitempty" db:"last_date"`
		Links          sql.NullString `json:"links,omitempty" db:"my_links"`
	}
	UserCreator struct {
		Id             int        `json:"-" db:"id"`
		Login          string     `json:"login" binding:"required"`
		EMail          string     `json:"email"`
		Password       string     `json:"password" db:"password_hash"`
		Confirmation   string     `json:"confirmation,omitempty" db:"confirmation"`
		IsConfirmation bool       `json:"is_confirm,omitempty" db:"is_confirm"`
		ExpiresAtToken TimeFormat `json:"expires_At_Token"`
		FirstName      string     `json:"firstName,omitempty" db:"first_name"`
		SecondName     string     `json:"secondName,omitempty"db:"second_name"`
		LastName       string     `json:"lastName,omitempty" db:"last_name"`
		ImagePath      string     `json:"imagePath,omitempty" db:"image"`
		Gender         int16      `json:"gender,omitempty" db:"gender"`
		Birthday       TimeFormat `json:"birthday,omitempty" db:"birthday"`
		Description    string     `json:"description,omitempty" db:"description"`
		RegisterDate   TimeFormat `json:"registerDate,omitempty" db:"register_date"`
		LastDate       TimeFormat `json:"lastDate,omitempty" db:"last_date"`
		Links          []Link     `json:"links,omitempty" db:"my_links"`
	}
	UserEmail struct {
		Id    int
		Email string
	}
	RefreshData struct {
		Token     string `json:"token" db:"refresh_token"`
		ExpiresAt int64  `json:"exp_token" db:"refresh_token_expires"`
	}
)

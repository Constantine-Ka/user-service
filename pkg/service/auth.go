package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/pkg/repository"
	"github.com/Constantine-Ka/user-service/tools/logging"
	"github.com/golang-jwt/jwt/v4"
	"net/mail"
	"strconv"
	"time"
)

const (
	salt         = "kkkkkkkkkkkkk"
	singingKey   = "thirteen13thirteenjwt"
	refreshKey   = "thirteen13thirteenrefresh"
	tokenTTL     = 168 * time.Hour  //7 days
	refreshTTL   = 4320 * time.Hour //6 Mouth
	refreshMouth = 6
)

type tokenClaims struct {
	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
	UserId    int   `json:"user_id"`
}
type refClaims struct {
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Subject   string `json:"subject"`
}

func (r refClaims) Valid() error {
	//TODO implement me
	panic("implement me")
}

func (t tokenClaims) Valid() error {
	//TODO implement me
	fmt.Println(t)
	return nil
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.UserCreator) (int, string, error) {

	logger := logging.GetLogger()
	logger.Info(user.Password)
	user.Password = generatePasswordHash(user.Password)

	user.Confirmation = generatePasswordHash(singingKey + user.Login)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (model.UserGet, string, error) {

	var (
		user model.UserGet
		err  error
	)

	if validEmail(username) {
		user, err = s.repo.GetUser("", username, generatePasswordHash(password), 0)
		if err != nil {
			return model.UserGet{}, "", err
		}
	} else {
		user, err = s.repo.GetUser(username, "", generatePasswordHash(password), 0)
		if err != nil {
			return model.UserGet{}, "", err
		}
	}
	if err != nil {
		return model.UserGet{}, "", err
	}
	user.Password = ""
	signingString, err := newJWT(user.Id)

	if err != nil {
		return model.UserGet{}, "", err
	}

	return user, signingString, err
}
func (s *AuthService) ConfirmUser(confirm string) (int, error) {
	logger := logging.GetLogger()
	var newToken model.RefreshData
	nowTime := time.Now().Unix()
	newToken.ExpiresAt = time.Now().Add(refreshTTL).Unix() //Six Mouth

	claims := jwt.NewWithClaims(jwt.SigningMethodHS384, &refClaims{
		ExpiresAt: newToken.ExpiresAt,
		IssuedAt:  nowTime,
		Subject:   confirm,
	})
	logger.Info(claims)
	signingString, err := claims.SignedString([]byte(refreshKey))
	logger.Info(signingString)
	logger.Error(err)
	if err != nil {
		return 0, err
	}

	newToken.Token = signingString
	return s.repo.ConfirmUser(confirm, newToken)
}
func (s *AuthService) ResetPassword(login, password string) error {
	var userdata model.UserCreator
	var num = isNumber(login)
	if validEmail(login) {
		userdata.EMail = login
	} else if num != 0 {
		userdata.Id = num
	} else if len(login) < 20 {
		userdata.Login = login
	} else {
		userdata.Confirmation = login
	}
	userdata.Password = generatePasswordHash(password)
	return s.repo.ResetPassword(userdata)
}
func (s *AuthService) GetConfirm(login string) (model.UserGet, error) {
	var column string
	var num = isNumber(login)
	if validEmail(login) {
		column = "email"
	} else if num != 0 {
		column = "id"
	} else {
		column = "login"
	}
	return s.repo.GetConfirm(login, column)
}
func (s *AuthService) RefreshToken(reftoken string, preId int) (int, string, error) {
	var (
		token string
		id    int
		err   error
	)

	id, err = s.repo.RefreshToken(reftoken)
	if err != nil {
		return 0, "", err
	}
	if preId != id {
		logging.GetLogger().Info(fmt.Sprintf("Attempt to refresh user-token id=%d with user refresh-token id=%d", preId, id))
		return 0, "", errors.New("invalid id")
	}
	token, err = newJWT(id)
	if err != nil {
		return 0, "", err
	}
	return id, token, err
}
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(singingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
func validEmail(login string) bool {
	_, err := mail.ParseAddress(login)
	return err == nil
}
func isNumber(id string) int {
	if num, err := strconv.Atoi(id); err == nil {
		return num
	}
	return 0
}
func newJWT(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(), //7 Day
		IssuedAt:  time.Now().Unix(),
		UserId:    id,
	})
	signingString, err := token.SignedString([]byte(singingKey))
	return signingString, err
}

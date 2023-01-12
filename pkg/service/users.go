package service

import (
	"encoding/json"
	"fmt"
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/pkg/repository"
	"github.com/Constantine-Ka/user-service/tools/logging"
)

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) MeUser(id int) (model.UserCreator, error) {
	//TODO implement me
	user, err := s.repo.GetUserOne(id)
	if err != nil {
		return model.UserCreator{}, err
	}
	user.Password = "it's Secret 😛😛😛"
	userNormal := converterDBtoNormal(user)
	return userNormal, nil
}
func (s *UserService) GetUserAll() ([]model.UserCreator, error) {
	var result []model.UserCreator
	temp, err := s.repo.GetUserAll()
	if err != nil {
		return nil, err
	}
	for _, get := range temp {
		fmt.Println(get)
		result = append(result, converterDBtoNormal(get))
	}
	return result, err
}
func (s *UserService) UpdateUser(id int, newData model.UserCreator) (model.UserCreator, error) {
	userGet, err := s.repo.UpdateUser(id, newData)
	user := converterDBtoNormal(userGet)
	return user, err
}
func (s *UserService) UpdateLink(arr []model.Link, id int) error {
	return s.repo.UpdateLink(arr, id)
}

func converterDBtoNormal(response model.UserGet) model.UserCreator {
	var result model.UserCreator
	{
		result.Id = response.Id
		result.Login = response.Login
		result.EMail = response.EMail
		result.Confirmation = response.Confirmation.String
		result.IsConfirmation = response.IsConfirmation.Bool
		result.ExpiresAtToken = model.TimeFormatter(response.ExpiresAt.Int64)
		result.FirstName = response.FirstName.String
		result.SecondName = response.SecondName.String
		result.LastName = response.LastName.String
		result.ImagePath = response.ImagePath.String
		result.Gender = response.Gender.Int16
		result.Birthday = model.TimeFormatter(response.Birthday.Int64)
		result.Description = response.Description.String
		result.RegisterDate = model.TimeFormatter(response.RegisterDate.Int64)
		result.LastDate = model.TimeFormatter(response.LastDate.Int64)
		//result.Links = response.Links.String
		err := json.Unmarshal([]byte(response.Links.String), &result.Links)
		if err != nil {
			logging.GetLogger().Info(err.Error())
			return result
		}
	}
	return result
}

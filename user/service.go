package user

import (
	"errors"
	"fmt"
	"ngepet-yuk/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAllUser() ([]UserFormat, error)
	SaveNewUser(user entity.UserInput) (UserFormat, error)
	GetUserByID(ID int) (UserFormat, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetAllUser() ([]UserFormat, error) {
	Users, err := s.repository.GetAll()
	var formatUsers []UserFormat

	for _, user := range Users {
		formatUser := FormatUser(user)
		formatUsers = append(formatUsers, formatUser)

	}
	if err != nil {
		return formatUsers, err
	}

	return formatUsers, nil
}

func (s *service) SaveNewUser(user entity.UserInput) (UserFormat, error) {
	//bcryp password jadi hash agar susah diretas
	genPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

	if err != nil {
		return UserFormat{}, err
	}

	var newUser = entity.User{
		UserName:  user.UserName,
		Email:     user.Email,
		Password:  string(genPassword), //diganti karena untuk hasing
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createUser, err := s.repository.Create(newUser)
	formatUser := FormatUser(createUser)

	if err != nil {
		return formatUser, err
	}

	return formatUser, nil
}

func (s *service) GetUserByID(ID int) (UserFormat, error) {
	user, err := s.repository.FindByID(ID)

	if err != nil {
		return UserFormat{}, err
	}

	if user.ID == 0 {
		newError := fmt.Sprintf("user id %d not found", user.ID)

		return UserFormat{}, errors.New(newError)
	}

	formatUser := FormatUser(user)

	return formatUser, nil
}

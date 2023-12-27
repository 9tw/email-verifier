package services

import (
	"email_verifier/features/user/domain"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &userService{qry: repo}
}

func (us *userService) Register(newUser domain.UserCore) (domain.UserCore, error) {
	generate, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error on bcrypt", err.Error())
		return domain.UserCore{}, errors.New("cannot encrypt password")
	}
	newUser.Password = string(generate)

	res, err := us.qry.AddUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}
	return res, nil
}

func (us *userService) Login(existUser domain.UserCore) (domain.UserCore, error) {
	res, err := us.qry.GetUser(existUser)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return domain.UserCore{}, errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("no data")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(existUser.Password))
	if err != nil {
		return domain.UserCore{}, errors.New("password not match")
	}

	return res, nil
}

func (us *userService) All() ([]domain.UserCore, error) {
	res, err := us.qry.GetAll()

	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("No Data")
		}
	}

	if len(res) == 0 {
		log.Info("No Data")
		return nil, errors.New("No Data")
	}
	return res, nil
}
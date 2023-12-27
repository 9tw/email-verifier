package repository

import (
	"email_verifier/features/user/domain"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type repoQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.Repository {
	return &repoQuery{db: db}
}

func (rq *repoQuery) AddUser(newUser domain.UserCore) (domain.UserCore, error) {
	var cnv User = FromDomain(newUser)
	if err := rq.db.Create(&cnv).Error; err != nil {
		log.Error("error on adding user", err.Error())
		return domain.UserCore{}, err
	}
	newUser = ToDomain(cnv)
	return newUser, nil
}

func (rq *repoQuery) GetUser(existUser domain.UserCore) (domain.UserCore, error) {
	var resQuery User
	if err := rq.db.First(&resQuery, "username = ?", existUser.Username).Error; err != nil {
		log.Error("error on get user login", err.Error())
		return domain.UserCore{}, nil
	}
	res := ToDomain(resQuery)
	return res, nil
}

func (rq *repoQuery) GetAll() ([]domain.UserCore, error) {
	var resQry []User

	if err := rq.db.Find(&resQry).Error; err != nil {
		return nil, err
	}

	res := ToDomainArray(resQry)
	return res, nil
}

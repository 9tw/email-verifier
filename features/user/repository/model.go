package repository

import (
	"email_verifier/features/user/domain"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName          string
	Username          string
	Password          string
	AppId             uint
	RoleId            uint
	UserType          string
	LastLogin         time.Time
	CreatedBy         int
	UpdatedBy         int
	Status            string
	EmailVerification string
	Timeout           time.Time
}

func FromDomain(du domain.UserCore) User {
	return User{
		Model:    gorm.Model{ID: du.ID},
		FullName: du.FullName,
		Username: du.Username,
		Password: du.Password,
		AppId:    du.AppId,
		RoleId:   du.RoleId,
		Status:   du.Status,
	}
}

func ToDomain(u User) domain.UserCore {
	return domain.UserCore{
		Model:    gorm.Model{ID: u.ID},
		FullName: u.FullName,
		Username: u.Username,
		Password: u.Password,
		AppId:    u.AppId,
		RoleId:   u.RoleId,
		Status:   u.Status,
	}
}

func ToDomainArray(au []User) []domain.UserCore {
	var res []domain.UserCore
	for _, val := range au {
		res = append(res, domain.UserCore{
			Model:    gorm.Model{ID: val.ID},
			FullName: val.FullName,
			Username: val.Username,
			AppId:    val.AppId,
			RoleId:   val.RoleId,
			Status:   val.Status,
		})
	}

	return res
}

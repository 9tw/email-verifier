package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserCore struct {
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

type Repository interface {
	AddUser(newUser UserCore) (UserCore, error)
	GetUser(existUser UserCore) (UserCore, error)
	GetSpesific(userID uint) (UserCore, error)
	PutActive(updatedUser UserCore, userID uint) (UserCore, error)
}

type Service interface {
	Register(newUser UserCore) (UserCore, error)
	Login(existUser UserCore) (UserCore, error)
	My(userID uint) (UserCore, error)
	Actived(updatedUser UserCore, userID uint) (UserCore, error)
}

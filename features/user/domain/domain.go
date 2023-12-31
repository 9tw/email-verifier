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
	GetUserWithUsername(username string) (UserCore, error)
	PutActive(updatedUser UserCore, user string) (UserCore, error)
}

type Service interface {
	Register(newUser UserCore) (UserCore, error)
	Login(existUser UserCore) (UserCore, error)
	My(userID uint) (UserCore, error)
	Actived(updatedUser UserCore, user string) (UserCore, error)
}

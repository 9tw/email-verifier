package services

import (
	"email_verifier/config"
	"email_verifier/features/user/domain"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

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
	newUser.Status = "0"

	// create random code for email
	// Go rune data type represent Unicode characters
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	// creat a random slice of runes (characters) to create our emailVerPassword (random string of characters)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	fmt.Println("emailVerRandRune:", emailVerRandRune)
	emailVerPassword := string(emailVerRandRune)
	fmt.Println("emailVerPassword:", emailVerPassword)
	var emailVerPWhash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	emailVerPWhash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err.Error())
		return domain.UserCore{}, errors.New("cannot encrypt email verifier")
	}
	fmt.Println("emailVerPWhash:", emailVerPWhash)
	newUser.EmailVerification = string(emailVerPWhash)
	// create u.timeout after 48 hours
	newUser.Timeout = time.Now().Local().AddDate(0, 0, 2)
	fmt.Println("timeout:", newUser.Timeout)

	res, err := us.qry.AddUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}
	subject := "Email Verificaion"
	HTMLbody :=
		`<html>
			<h1>Click Link to Veify Email</h1>
			<a href="http://localhost:8000/emailver/` + newUser.Username + `/` + emailVerPassword + `">click to verify email</a>
		</html>`
	fmt.Println(HTMLbody)
	err = config.SendEmail(subject, HTMLbody, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "error") {
			return domain.UserCore{}, errors.New("can't send email")
		}
		return domain.UserCore{}, errors.New("some problem on sending email")
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

func (us *userService) My(userID uint) (domain.UserCore, error) {
	res, err := us.qry.GetSpesific(userID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return domain.UserCore{}, errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("no data")
		}
	}
	return res, nil
}

func (us *userService) Actived(updatedUser domain.UserCore, userID uint) (domain.UserCore, error) {
	updatedUser.Status = "1"
	res, err := us.qry.PutActive(updatedUser, userID)
	if err != nil {
		if strings.Contains(err.Error(), "column") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}

	return res, nil
}

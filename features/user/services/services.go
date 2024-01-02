package services

import (
	"bufio"
	"email_verifier/config"
	"email_verifier/features/user/domain"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/labstack/gommon/log"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

var (
	verifier = emailverifier.NewVerifier()
)

func init() {
	verifier = verifier.EnableDomainSuggest()
	dispEmailsDomains := MustDispEmailDom()
	verifier = verifier.AddDisposableDomains(dispEmailsDomains)
}

func MustDispEmailDom() (dispEmailDomains []string) {
	file, err := os.Open("./disposable_email_blocklist.txt")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dispEmailDomains = append(dispEmailDomains, scanner.Text())
	}
	return dispEmailDomains
}

type userService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &userService{qry: repo}
}

func (us *userService) Register(newUser domain.UserCore) (domain.UserCore, error) {
	// check username for only alphaNumeric characters
	for _, char := range newUser.FullName {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return domain.UserCore{}, errors.New("only alphanumeric characters allowed for username")
		}
	}

	// check username length
	if 5 >= len(newUser.FullName) && len(newUser.FullName) >= 50 {
		return domain.UserCore{}, errors.New("username length must be greater than 4 and less than 51 characters")
	}

	err := passwordvalidator.Validate(newUser.Password, minEntropyBits)
	if err != nil {
		return domain.UserCore{}, err
	}

	_, err = us.qry.GetUserWithUsername(newUser.Username)
	if err != nil {
		return domain.UserCore{}, errors.New("Username already exists")
	}

	mail, err := verifier.Verify(newUser.Username)
	if err != nil {
		return domain.UserCore{}, err
	}

	// check syntax, needs @ and . for starters
	if !mail.Syntax.Valid {
		return domain.UserCore{}, errors.New("Email address syntax is invalid")
	}

	// check if disposable
	if mail.Disposable {
		return domain.UserCore{}, errors.New("We do not accept disposable email addresses")
	}

	// check if there is domain Suggestion
	if mail.Suggestion != "" {
		return domain.UserCore{}, errors.New("Email address is not reachable")
	}

	// possible return string values: yes, no, unkown
	if mail.Reachable == "no" {
		return domain.UserCore{}, errors.New("Email address is not reachable")
	}

	// check MX records so we know DNS setup properly to recieve emails
	if !mail.HasMxRecords {
		return domain.UserCore{}, errors.New("Domain entered not properly setup to recieve emails, MX record not found")
	}

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

func (us *userService) Actived(updatedUser domain.UserCore, user string) (domain.UserCore, error) {
	updatedUser.Status = "1"
	res, err := us.qry.PutActive(updatedUser, user)
	if err != nil {
		if strings.Contains(err.Error(), "column") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}

	return res, nil
}

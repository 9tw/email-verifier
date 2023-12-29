package delivery

import "email_verifier/features/user/domain"

type UserFormat struct {
	FullName string `json:"fullname" form:"fullname"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type LoginFormat struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type ActiveFormat struct {
	Username          string `json:"username" form:"username"`
	EmailVerification string `json:"verification" form:"verification"`
}

func ToDomain(i interface{}) domain.UserCore {
	switch i.(type) {
	case UserFormat:
		cnv := i.(UserFormat)
		return domain.UserCore{
			FullName: cnv.FullName,
			Username: cnv.Username,
			Password: cnv.Password,
		}
	case LoginFormat:
		cnv := i.(LoginFormat)
		return domain.UserCore{
			Username: cnv.Username,
			Password: cnv.Password,
		}
	case ActiveFormat:
		cnv := i.(ActiveFormat)
		return domain.UserCore{
			Username:          cnv.Username,
			EmailVerification: cnv.EmailVerification,
		}
	}
	return domain.UserCore{}
}

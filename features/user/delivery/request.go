package delivery

import "email_verifier/features/user/domain"

type UserFormat struct {
	FullName string `json:"fullname" form:"fullname"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	AppId    uint   `json:"app_id" form:"app_id"`
	RoleId   uint   `json:"role_id" form:"role_id"`
	Status   string `json:"status" form:"status"`
}

type LoginFormat struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func ToDomain(i interface{}) domain.UserCore {
	switch i.(type) {
	case UserFormat:
		cnv := i.(UserFormat)
		return domain.UserCore{FullName: cnv.FullName, Username: cnv.Username, Password: cnv.Password, AppId: cnv.AppId, RoleId: cnv.RoleId, Status: cnv.Status}
	case LoginFormat:
		cnv := i.(LoginFormat)
		return domain.UserCore{Username: cnv.Username, Password: cnv.Password}
	}
	return domain.UserCore{}
}

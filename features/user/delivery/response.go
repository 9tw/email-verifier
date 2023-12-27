package delivery

import "email_verifier/features/user/domain"

type RegisterResponse struct {
	FullName string `json:"name"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Username string `json:"username"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "user":
		cnv := core.(domain.UserCore)
		res = RegisterResponse{FullName: cnv.FullName, Username: cnv.Username}
	case "login":
		cnv := core.(domain.UserCore)
		res = LoginResponse{Username: cnv.Username}
	case "all":
		var arr []RegisterResponse
		cnv := core.([]domain.UserCore)
		for _, val := range cnv {
			arr = append(arr, RegisterResponse{FullName: val.FullName, Username: val.Username})
		}
		res = arr
	}
	return res
}

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Message": msg,
		"Data":    data,
	}
}

func SuccessLogin(msg string, data interface{}, token interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Message": msg,
		"Data":    data,
		"Token":   token,
	}
}

func FailResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"Message": msg,
	}
}

package delivery

import (
	"email_verifier/config"
	"email_verifier/features/user/domain"
	"email_verifier/utils/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type userHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := userHandler{srv: srv}
	e.POST("/register", handler.Register())
	e.POST("/login", handler.Login())
	e.GET("/emailver/:username/:verPass", handler.EmailVerifier())
	e.GET("/user", handler.MyProfile(), middleware.JWT([]byte(config.JwtKey)))
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input UserFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		cnv := ToDomain(input)
		res, err := uh.srv.Register(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}
		return c.JSON(http.StatusCreated, SuccessResponse("Success create new user", ToResponse(res, "user")))
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		cnv := ToDomain(input)
		res, err := uh.srv.Login(cnv)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		tkn := common.GenerateToken(uint(res.ID))

		return c.JSON(http.StatusAccepted, SuccessLogin("Success to login", ToResponse(res, "login"), tkn))
	}
}

func (uh *userHandler) EmailVerifier() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user = c.Param("username")
		var link = c.Param("verPass")
		if user == "" {
			return c.JSON(http.StatusUnauthorized, FailResponse("User doesn't exists."))
		} else {
			var input ActiveFormat
			if err := c.Bind(&input); err != nil {
				return c.JSON(http.StatusBadRequest, FailResponse("cannot bind update data"))
			}

			cnv := ToDomain(input)
			cnv.Username = user
			cnv.EmailVerification = link
			res, err := uh.srv.Actived(cnv, user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get my profile.", ToResponse(res, "user")))
		}
	}
}

func (uh *userHandler) MyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := common.ExtractToken(c)
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, FailResponse("ID doesn't exists."))
		} else {
			res, err := uh.srv.My(uint(userID))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
			}
			return c.JSON(http.StatusOK, SuccessResponse("Success get my profile.", ToResponse(res, "user")))
		}
	}
}

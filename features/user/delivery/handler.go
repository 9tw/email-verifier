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
	e.GET("/users", handler.AllAccounts(), middleware.JWT([]byte(config.JwtKey)))
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

func (uh *userHandler) AllAccounts() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := uh.srv.All()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, SuccessResponse("Success get all accounts.", ToResponse(res, "all")))
	}
}

package middleware

import (
	"charum/business/users"
	"charum/helper"
	"charum/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoleMiddleware struct {
	Role           []string
	UserRepository users.Repository
	Func           echo.HandlerFunc
}

func (rm RoleMiddleware) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := util.GetUIDFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid token",
				Data:    nil,
			})
		}

		user, err := rm.UserRepository.GetByID(uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
		}

		for _, role := range rm.Role {
			if user.Role == role {
				return next(c)
			}
		}

		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    nil,
		})
	}
}

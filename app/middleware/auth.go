package middleware

import (
	"charum/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type RoleMiddleware struct {
	Role []string
	Func echo.HandlerFunc
}

func (rm RoleMiddleware) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		token := strings.Replace(authHeader, "Bearer ", "", -1)

		claims, err := util.GetPayloadToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for _, role := range rm.Role {
			if claims.Role == role {
				return next(c)
			}
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
	}
}

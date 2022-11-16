package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	apiV1 := e.Group("/api/v1")

	apiV1.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World!",
		})
	})
}

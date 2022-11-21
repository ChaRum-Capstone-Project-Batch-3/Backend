package route

import (
	_middleware "charum/app/middleware"
	"charum/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware echo.MiddlewareFunc
	UserController   *users.UserController
}

func (cl *ControllerList) Init(e *echo.Echo) {
	_middleware.InitLogger(e)

	userMiddleware := _middleware.RoleMiddleware{Role: []string{"user"}}
	adminMiddleware := _middleware.RoleMiddleware{Role: []string{"admin"}}

	apiV1 := e.Group("/api/v1")
	apiV1.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World!",
		})
	})

	user := apiV1.Group("/user")
	user.POST("/register", cl.UserController.UserRegister)
	user.POST("/login", cl.UserController.Login)
	user.GET("/user", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello User!",
		})
	}, userMiddleware.Check)
	user.GET("/admin", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello Admin!",
		})
	}, adminMiddleware.Check)

	admin := apiV1.Group("/admin")
	admin.GET("/user/:page", cl.UserController.GetAllUser, adminMiddleware.Check)
}

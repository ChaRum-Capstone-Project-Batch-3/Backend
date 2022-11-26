package route

import (
	_middleware "charum/app/middleware"
	"charum/controller/topics"
	"charum/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware echo.MiddlewareFunc
	UserController   *users.UserController
	TopicController  *topics.TopicController
}

func (cl *ControllerList) Init(e *echo.Echo) {
	_middleware.InitLogger(e)

	// userMiddleware := _middleware.RoleMiddleware{Role: []string{"user"}}
	adminMiddleware := _middleware.RoleMiddleware{Role: []string{"admin"}}
	authMiddleware := _middleware.RoleMiddleware{Role: []string{"user", "admin"}}

	apiV1 := e.Group("/api/v1")
	apiV1.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to Charum!",
		})
	})

	user := apiV1.Group("/user")
	user.POST("/register", cl.UserController.Register)
	user.POST("/login", cl.UserController.Login)
	user.GET("/profile", cl.UserController.GetProfile, authMiddleware.Check)

	admin := apiV1.Group("/admin", adminMiddleware.Check)

	adminUser := admin.Group("/user")
	adminUser.GET("/:page", cl.UserController.GetManyWithPagination)
	adminUser.GET("/id/:id", cl.UserController.GetByID)
	adminUser.PUT("/id/:id", cl.UserController.Update)
	adminUser.DELETE("/id/:id", cl.UserController.Delete)
	adminUser.PUT("/suspend/:id", cl.UserController.Suspend)
	adminUser.PUT("/unsuspend/:id", cl.UserController.Unsuspend)

	adminTopic := admin.Group("/topic")
	adminTopic.POST("/create", cl.TopicController.CreateTopic)
	adminTopic.GET("/id/:id", cl.TopicController.GetByID)
	adminTopic.PUT("/id/:id", cl.TopicController.UpdateTopic)
	adminTopic.DELETE("/id/:id", cl.TopicController.DeleteTopic)
}

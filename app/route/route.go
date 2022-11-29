package route

import (
	_middleware "charum/app/middleware"
	"charum/controller/bookmarks"
	"charum/controller/threads"
	"charum/controller/topics"
	"charum/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware   echo.MiddlewareFunc
	UserController     *users.UserController
	TopicController    *topics.TopicController
	ThreadController   *threads.ThreadController
	BookmarkController *bookmarks.BookmarkController
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
	user.POST("/bookmark/:thread_id", cl.BookmarkController.AddBookmark, authMiddleware.Check)

	thread := apiV1.Group("/thread")
	thread.POST("", cl.ThreadController.Create, authMiddleware.Check)
	thread.GET("/:page", cl.ThreadController.GetManyWithPagination)
	thread.GET("/id/:id", cl.ThreadController.GetByID)
	thread.PUT("/id/:id", cl.ThreadController.Update, authMiddleware.Check)
	thread.DELETE("/id/:id", cl.ThreadController.Delete, authMiddleware.Check)

	admin := apiV1.Group("/admin", adminMiddleware.Check)

	adminUser := admin.Group("/user")
	adminUser.GET("/:page", cl.UserController.GetManyWithPagination)
	adminUser.GET("/id/:id", cl.UserController.GetByID)
	adminUser.PUT("/id/:id", cl.UserController.Update)
	adminUser.DELETE("/id/:id", cl.UserController.Delete)
	adminUser.PUT("/suspend/:id", cl.UserController.Suspend)
	adminUser.PUT("/unsuspend/:id", cl.UserController.Unsuspend)

	adminTopic := admin.Group("/topic")
	adminTopic.POST("", cl.TopicController.CreateTopic)
	adminTopic.GET("", cl.TopicController.GetAll)
	adminTopic.GET("/name/:topic", cl.TopicController.GetByTopic)
	adminTopic.GET("/:id", cl.TopicController.GetByID)
	adminTopic.PUT("/:id", cl.TopicController.UpdateTopic)
	adminTopic.DELETE("/:id", cl.TopicController.DeleteTopic)
}

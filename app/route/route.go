package route

import (
	_middleware "charum/app/middleware"
	_usersDomain "charum/business/users"
	"charum/controller/comments"
	followThreads "charum/controller/follow_threads"
	"charum/controller/threads"
	"charum/controller/topics"
	"charum/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware       echo.MiddlewareFunc
	UserRepository         _usersDomain.Repository
	UserController         *users.UserController
	TopicController        *topics.TopicController
	ThreadController       *threads.ThreadController
	CommentController      *comments.CommentController
	FollowThreadController *followThreads.FollowThreadController
}

func (cl *ControllerList) Init(e *echo.Echo) {
	_middleware.InitLogger(e)

	adminMiddleware := _middleware.RoleMiddleware{Role: []string{"admin"}, UserRepository: cl.UserRepository}
	authMiddleware := _middleware.RoleMiddleware{Role: []string{"user", "admin"}, UserRepository: cl.UserRepository}

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
	user.PUT("/profile", cl.UserController.UserUpdate, authMiddleware.Check)

	thread := apiV1.Group("/thread")
	thread.POST("", cl.ThreadController.Create, authMiddleware.Check)
	thread.GET("/:page", cl.ThreadController.GetManyWithPagination)

	threadID := thread.Group("/id")
	threadID.GET("/:id", cl.ThreadController.GetByID)
	threadID.PUT("/:id", cl.ThreadController.Update, authMiddleware.Check)
	threadID.DELETE("/:id", cl.ThreadController.Delete, authMiddleware.Check)

	threadFollow := thread.Group("/follow")
	threadFollow.GET("", cl.FollowThreadController.GetFollowedThreadByToken, authMiddleware.Check)
	threadFollow.GET("/:user-id", cl.FollowThreadController.GetFollowedThreadByUserID)
	threadFollow.POST("/:thread-id", cl.FollowThreadController.Create, authMiddleware.Check)
	threadFollow.DELETE("/:thread-id", cl.FollowThreadController.Delete, authMiddleware.Check)

	threadComment := thread.Group("/comment")
	threadComment.POST("/:thread-id", cl.CommentController.Create, authMiddleware.Check)
	threadComment.PUT("/:comment-id", cl.CommentController.Update, authMiddleware.Check)
	threadComment.DELETE("/:comment-id", cl.CommentController.Delete, authMiddleware.Check)

	admin := apiV1.Group("/admin", adminMiddleware.Check)

	adminUser := admin.Group("/user")
	adminUser.GET("/:page", cl.UserController.GetManyWithPagination)
	adminUser.PUT("/suspend/:id", cl.UserController.Suspend)
	adminUser.PUT("/unsuspend/:id", cl.UserController.Unsuspend)

	adminUserID := adminUser.Group("/id")
	adminUserID.GET("/:id", cl.UserController.GetByID)
	adminUserID.PUT("/:id", cl.UserController.AdminUpdate)
	adminUserID.DELETE("/:id", cl.UserController.Delete)

	adminTopic := admin.Group("/topic")
	adminTopic.POST("", cl.TopicController.CreateTopic)
	adminTopic.GET("", cl.TopicController.GetAll)
	adminTopic.GET("/:id", cl.TopicController.GetByID)
	adminTopic.PUT("/:id", cl.TopicController.UpdateTopic)
	adminTopic.DELETE("/:id", cl.TopicController.DeleteTopic)
}

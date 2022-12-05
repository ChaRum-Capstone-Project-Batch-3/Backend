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

	topic := apiV1.Group("/topic")
	topic.GET("", cl.TopicController.GetAll)
	topic.GET("/:topic-id", cl.TopicController.GetByID)

	thread := apiV1.Group("/thread")
	thread.GET("", cl.ThreadController.GetManyByToken, authMiddleware.Check)
	thread.POST("", cl.ThreadController.Create, authMiddleware.Check)
	thread.POST("/:page", cl.ThreadController.GetManyWithPagination)

	threadID := thread.Group("/id")
	threadID.GET("/:thread-id", cl.ThreadController.GetByID)
	threadID.PUT("/:thread-id", cl.ThreadController.UserUpdate, authMiddleware.Check)
	threadID.DELETE("/:thread-id", cl.ThreadController.UserDelete, authMiddleware.Check)

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
	adminUser.POST("/:page", cl.UserController.GetManyWithPagination)
	adminUser.PUT("/suspend/:user-id", cl.UserController.Suspend)
	adminUser.PUT("/unsuspend/:user-id", cl.UserController.Unsuspend)

	adminUserID := adminUser.Group("/id")
	adminUserID.GET("/:user-id", cl.UserController.GetByID)
	adminUserID.PUT("/:user-id", cl.UserController.AdminUpdate)
	adminUserID.DELETE("/:user-id", cl.UserController.Delete)

	adminTopic := admin.Group("/topic")
	adminTopic.POST("", cl.TopicController.Create)
	adminTopic.GET("", cl.TopicController.GetAll)
	adminTopic.GET("/:topic-id", cl.TopicController.GetByID)
	adminTopic.PUT("/:topic-id", cl.TopicController.Update)
	adminTopic.DELETE("/:topic-id", cl.TopicController.Delete)

	adminThread := admin.Group("/thread")
	adminThread.POST("/:page", cl.ThreadController.GetManyWithPagination)
	adminThread.GET("/id/:thread-id", cl.ThreadController.GetByID)
	adminThread.PUT("/id/:thread-id", cl.ThreadController.AdminUpdate)
	adminThread.DELETE("/id/:thread-id", cl.ThreadController.AdminDelete)
}

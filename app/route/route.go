package route

import (
	_middleware "charum/app/middleware"
	_usersDomain "charum/business/users"
	_bookmarkController "charum/controller/bookmarks"
	"charum/controller/comments"
	followThreads "charum/controller/follow_threads"
	"charum/controller/forgot_password"
	"charum/controller/reports"
	"charum/controller/threads"
	"charum/controller/topics"
	"charum/controller/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ControllerList struct {
	LoggerMiddleware         echo.MiddlewareFunc
	UserRepository           _usersDomain.Repository
	UserController           *users.UserController
	TopicController          *topics.TopicController
	ThreadController         *threads.ThreadController
	CommentController        *comments.CommentController
	FollowThreadController   *followThreads.FollowThreadController
	BookmarkController       *_bookmarkController.BookmarkController
	ForgotPasswordController *forgot_password.ForgotPasswordController
	ReportController         *reports.ReportController
}

func (cl *ControllerList) Init(e *echo.Echo) {
	_middleware.InitLogger(e)

	e.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to Charum!",
		})
	})

	apiV1 := e.Group("/api/v1")

	user := apiV1.Group("/user")
	user.POST("/register", cl.UserController.Register)
	user.POST("/login", cl.UserController.Login)
	user.GET("/profile", cl.UserController.GetProfile, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	user.PUT("/profile", cl.UserController.UserUpdate, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	user.POST("/report/:user-id", cl.ReportController.ReportUser, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	user.PUT("/change-password", cl.UserController.UpdatePassword, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	user.POST("/forgot-password", cl.ForgotPasswordController.Generate)
	user.GET("/forgot-password/:token", cl.ForgotPasswordController.ValidateToken)
	user.POST("/forgot-password/:token", cl.ForgotPasswordController.Update)

	topic := apiV1.Group("/topic")
	topic.GET("/:page", cl.TopicController.GetManyWithPagination)
	topic.GET("/id/:topic-id", cl.TopicController.GetByID)

	thread := apiV1.Group("/thread")
	thread.GET("", cl.ThreadController.GetManyByToken, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	thread.POST("", cl.ThreadController.Create, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	thread.GET("/:page", cl.ThreadController.GetManyWithPagination)
	threadID := thread.Group("/id")
	threadID.GET("/:thread-id", cl.ThreadController.GetByID)
	threadID.PUT("/:thread-id", cl.ThreadController.UserUpdate, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadID.DELETE("/:thread-id", cl.ThreadController.UserDelete, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadFollow := thread.Group("/follow")
	threadFollow.GET("", cl.FollowThreadController.GetFollowedThreadByToken, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadFollow.GET("/:user-id", cl.FollowThreadController.GetFollowedThreadByUserID)
	threadFollow.POST("/:thread-id", cl.FollowThreadController.Create, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadFollow.DELETE("/:thread-id", cl.FollowThreadController.Delete, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadComment := thread.Group("/comment")
	threadComment.POST("/:thread-id", cl.CommentController.Create, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadComment.PUT("/:comment-id", cl.CommentController.Update, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadComment.DELETE("/:comment-id", cl.CommentController.Delete, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadLike := thread.Group("/like")
	threadLike.GET("", cl.ThreadController.GetLikedThreadByToken, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadLike.GET("/:user-id", cl.ThreadController.GetLikedThreadByUserID)
	threadLike.POST("/id/:thread-id", cl.ThreadController.Like, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadLike.DELETE("/id/:thread-id", cl.ThreadController.Unlike, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadBookmark := thread.Group("/bookmark")
	threadBookmark.GET("", cl.BookmarkController.GetAllByToken, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadBookmark.POST("/:thread-id", cl.BookmarkController.Create, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadBookmark.DELETE("/:thread-id", cl.BookmarkController.Delete, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))
	threadReport := thread.Group("/report")
	threadReport.POST("/:thread-id", cl.ReportController.ReportThread, _middleware.Check([]string{"user", "admin"}, cl.UserRepository))

	// Admin
	admin := apiV1.Group("/admin", _middleware.Check([]string{"admin"}, cl.UserRepository))
	admin.GET("/statistics", cl.ReportController.CountAllData)

	adminUser := admin.Group("/user")
	adminUser.GET("/:page", cl.UserController.GetManyWithPagination)
	adminUser.PUT("/suspend/:user-id", cl.UserController.Suspend)
	adminUser.PUT("/unsuspend/:user-id", cl.UserController.Unsuspend)
	adminUser.GET("/report", cl.ReportController.GetAllReportedUsers)
	adminUser.GET("/report/:user-id", cl.ReportController.GetUserReportedID)
	adminUserID := adminUser.Group("/id")
	adminUserID.GET("/:user-id", cl.UserController.GetByID)
	adminUserID.PUT("/:user-id", cl.UserController.AdminUpdate)
	adminUserID.DELETE("/:user-id", cl.UserController.Delete)

	adminReport := admin.Group("/report")
	adminReport.GET("", cl.ReportController.GetAll)

	adminTopic := admin.Group("/topic")
	adminTopic.POST("", cl.TopicController.Create)
	adminTopic.GET("/:page", cl.TopicController.GetManyWithPagination)
	adminTopic.GET("/id/:topic-id", cl.TopicController.GetByID)
	adminTopic.PUT("/id/:topic-id", cl.TopicController.Update)
	adminTopic.DELETE("/id/:topic-id", cl.TopicController.Delete)

	adminThread := admin.Group("/thread")
	adminThread.GET("/:page", cl.ThreadController.GetManyWithPagination)
	adminThread.GET("/report", cl.ReportController.GetAllReportedThreads)
	adminThread.GET("/report/:thread-id", cl.ReportController.GetThreadReportedID)
	adminThread.GET("/id/:thread-id", cl.ThreadController.GetByID)
	adminThread.PUT("/id/:thread-id", cl.ThreadController.AdminUpdate)
	adminThread.DELETE("/id/:thread-id", cl.ThreadController.AdminDelete)

}

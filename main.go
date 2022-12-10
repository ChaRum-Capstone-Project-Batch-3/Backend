package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_route "charum/app/route"
	_driver "charum/driver"
	_mongo "charum/driver/mongo"
	_util "charum/util"

	_userUseCase "charum/business/users"
	_userController "charum/controller/users"

	_topicUseCase "charum/business/topics"
	_topicController "charum/controller/topics"

	_threadUseCase "charum/business/threads"
	_threadController "charum/controller/threads"

	_commentUseCase "charum/business/comments"
	_commentController "charum/controller/comments"

	_followThreadUseCase "charum/business/follow_threads"
	_followThreadController "charum/controller/follow_threads"

	_bookmarkUseCase "charum/business/bookmarks"
	_bookmarkController "charum/controller/bookmarks"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))
	cloudinary := _driver.NewCloudinaryRepository()

	userRepository := _driver.NewUserRepository(database)
	topicRepository := _driver.NewTopicRepository(database)
	threadRepository := _driver.NewThreadRepository(database)
	commentRepository := _driver.NewCommentRepository(database)
	followThreadRepository := _driver.NewFollowThreadRepository(database)
	bookmarkRepository := _driver.NewBookmarkRepository(database)

	userUsecase := _userUseCase.NewUserUseCase(userRepository, cloudinary)
	topicUsecase := _topicUseCase.NewTopicUseCase(topicRepository, cloudinary)
	threadUsecase := _threadUseCase.NewThreadUseCase(threadRepository, topicRepository, userRepository)
	commentUsecase := _commentUseCase.NewCommentUseCase(commentRepository, threadRepository, userRepository)
	followThreadUsecase := _followThreadUseCase.NewFollowThreadUseCase(followThreadRepository, userRepository, threadRepository, commentRepository, threadUsecase)
	bookmarkUsecase := _bookmarkUseCase.NewBookmarkUseCase(bookmarkRepository, threadRepository, userRepository, topicRepository, threadUsecase)

	userController := _userController.NewUserController(userUsecase, threadUsecase, commentUsecase, followThreadUsecase, bookmarkUsecase)
	topicController := _topicController.NewTopicController(topicUsecase, threadUsecase, commentUsecase, followThreadUsecase, bookmarkUsecase)
	threadController := _threadController.NewThreadController(threadUsecase, commentUsecase, followThreadUsecase, userUsecase, bookmarkUsecase)
	commentController := _commentController.NewCommentController(commentUsecase, followThreadUsecase)
	followThreadController := _followThreadController.NewFollowThreadController(followThreadUsecase)
	bookmarkController := _bookmarkController.NewBookmarkController(bookmarkUsecase, followThreadUsecase, commentUsecase)

	routeController := _route.ControllerList{
		UserRepository:         userRepository,
		UserController:         userController,
		TopicController:        topicController,
		ThreadController:       threadController,
		CommentController:      commentController,
		FollowThreadController: followThreadController,
		BookmarkController:     bookmarkController,
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	routeController.Init(e)

	appPort := fmt.Sprintf(":%s", _util.GetConfig("APP_PORT"))

	go func() {
		if err := e.StartTLS(appPort, "cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	wait := _util.GracefulShutdown(context.Background(), 2*time.Second, map[string]_util.Operation{
		"database": func(ctx context.Context) error {
			return _mongo.Close(database)
		},
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(context.Background())
		},
	})

	<-wait
}

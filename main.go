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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))

	userRepository := _driver.NewUserRepository(database)
	userUsecase := _userUseCase.NewUserUseCase(userRepository)
	userController := _userController.NewUserController(userUsecase)

	topicRepository := _driver.NewTopicRepository(database)
	topicUsecase := _topicUseCase.NewTopicUseCase(topicRepository)
	topicController := _topicController.NewTopicController(topicUsecase)

	threadRepository := _driver.NewThreadRepository(database)
	threadUsecase := _threadUseCase.NewThreadUseCase(threadRepository, topicRepository, userRepository)
	threadController := _threadController.NewThreadController(threadUsecase)

	commentRepository := _driver.NewCommentRepository(database)
	commentUsecase := _commentUseCase.NewCommentUseCase(commentRepository, threadRepository)
	commentController := _commentController.NewCommentController(commentUsecase)

	routeController := _route.ControllerList{
		UserController:    userController,
		TopicController:   topicController,
		ThreadController:  threadController,
		CommentController: commentController,
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
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

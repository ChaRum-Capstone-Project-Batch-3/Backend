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

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))

	userRepository := _driver.NewUserRepository(database)
	userUsecase := _userUseCase.NewUserUseCase(userRepository)
	userController := _userController.NewUserController(userUsecase)

	routeController := _route.ControllerList{
		UserController: userController,
	}

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

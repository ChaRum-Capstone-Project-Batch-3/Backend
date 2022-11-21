package main

import (
	"fmt"

	_route "charum/app/route"
	_driver "charum/driver"
	_mongo "charum/driver/mongo"
	_util "charum/util"

	_userUseCase "charum/businesses/users"
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
	e.Logger.Fatal(e.Start(appPort))
}

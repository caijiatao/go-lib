package main

import (
	"github.com/gin-gonic/gin"
	"golib/libs/gin_helper"
	"golib/libs/logger"
	"sync"
)

type Application struct {
	Engine     *gin.Engine
	Controller *ApplicationController
}

type ApplicationController struct {
}

var application *Application
var appOnce sync.Once

func App() *Application {
	appOnce.Do(func() {
		application = &Application{
			Controller: &ApplicationController{},
		}
	})
	return application
}

func (app *Application) Init(engine *gin.Engine) *Application {
	app.Engine = engine
	app.InitMiddleware().InitRouter()
	return app
}

func (app *Application) InitRouter() *Application {
	gin_helper.RegisterAllRoutes(app.Controller, app.Engine)
	return app
}

func (app *Application) InitMiddleware() *Application {
	app.Engine.Use(logger.MiddlewareLoggerInfo())
	app.Engine.Use(gin_helper.ErrorHandler())
	app.Engine.Use(gin_helper.Context())
	return app
}

func (app *Application) Run() error {
	gin_helper.Init()
	serverConfig := gin_helper.GetServerConfig()
	return app.Engine.Run(":" + serverConfig.HttpPort)
}

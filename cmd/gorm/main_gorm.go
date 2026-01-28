package main

import (
	"sdm_demo_todolist/internal/gorm/dbal"
	handlers2 "sdm_demo_todolist/internal/gorm/handlers"
	"sdm_demo_todolist/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	err := dbal.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = dbal.CloseDB()
	}()

	gin.SetMode(gin.ReleaseMode)

	myRouter := gin.New()

	projectHandlers := handlers2.NewGormProjectHandlers()
	taskHandlers := handlers2.NewGormTaskHandlers()

	etc.AssignHandlers(myRouter, "Go, Gorm, SQLite3", projectHandlers, taskHandlers)

	etc.Listen(myRouter)
}

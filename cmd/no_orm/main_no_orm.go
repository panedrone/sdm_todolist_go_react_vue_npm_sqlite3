package main

import (
	"sdm_demo_todolist/internal/no_orm/dbal"
	handlers2 "sdm_demo_todolist/internal/no_orm/handlers"
	"sdm_demo_todolist/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	err := dbal.OpenDB()
	if err != nil {
		println(err.Error())
		return
	}
	defer func() {
		_ = dbal.CloseDB()
	}()

	gin.SetMode(gin.ReleaseMode)

	myRouter := gin.New()

	projectHandlers := handlers2.NewNoOrmProjectHandlers()
	taskHandlers := handlers2.NewNoOrmTaskHandlers()

	etc.AssignHandlers(myRouter, "Go, database/sql, SQLite3", projectHandlers, taskHandlers)

	etc.Listen(myRouter)
}

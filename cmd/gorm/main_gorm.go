package main

import (
	"sdm_demo_todolist/internal/gorm/dbal"
	"sdm_demo_todolist/internal/gorm/handlers"
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

	//myRouter.Use(WithDS())

	projectHandlers := handlers.NewGormProjectHandlers()
	taskHandlers := handlers.NewGormTaskHandlers()

	etc.AssignHandlers(myRouter, "Go, Gorm, SQLite3", projectHandlers, taskHandlers)

	etc.Listen(myRouter)
}

//func WithDS() gin.HandlerFunc {
//	const dsKey = "DataStore"
//	dbal.DS = func(ctx context.Context) dbal.DataStore {
//		g := ctx.(*gin.Context)
//		return g.MustGet(dsKey).(dbal.DataStore)
//	}
//	return func(c *gin.Context) {
//		ds := dbal.GetInstance()
//		c.Set(dsKey, ds)
//		c.Next()
//	}
//}

//func WithDS() gin.HandlerFunc {
//	dbal.DS = func(ctx context.Context) dbal.DataStore {
//		return dbal.GetInstance()
//	}
//	return func(c *gin.Context) {
//		c.Next()
//	}
//}

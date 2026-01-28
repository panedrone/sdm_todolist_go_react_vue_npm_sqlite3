package main

import (
	"fmt"
	"os"
	"runtime"
	"sdm_demo_todolist/internal/sqlx/dbal"
	handlers2 "sdm_demo_todolist/internal/sqlx/handlers"
	"sdm_demo_todolist/internal/sqlx/swagger"
	"sdm_demo_todolist/pkg"

	"github.com/gin-gonic/gin"
)

// @schemes	http
// @produce	json
// @version	0.0.1
//
// @title		SDM TodoApp API
// @in			header
// @BasePath	/api
// @accept		json
// @host		127.0.0.1:8080
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

	swagger.Init(myRouter)

	myOS, myArch := runtime.GOOS, runtime.GOARCH
	inContainer := " docker,"
	if _, err := os.Lstat("/.dockerenv"); err != nil && os.IsNotExist(err) {
		inContainer = ""
	}

	// whoIam := fmt.Sprintf(`%v, %v,%v sqlx, sqlite3, <a target="_blank" href="swagger/index.html">swagger</a>`, myOS, myArch, inContainer)
	whoIam := fmt.Sprintf(`%v, %v,%v sqlx, sqlite3`, myOS, myArch, inContainer)

	etc.AssignHandlers(myRouter, whoIam, handlers2.NewSqlxProjectHandlers(), handlers2.NewSqlxTaskHandlers())

	etc.Listen(myRouter)
}

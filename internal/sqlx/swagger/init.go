package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	////////////////////////////////////////////////////////////
	// https://github.com/swaggo/swag/issues/1019
	// You need to import docs package in your main package.
	// docs package is created by swagger when you swagger init.
	// === panedrone: !!! import _ "sdm_demo_todolist/sqlx/docs"
	////////////////////////////////////////////////////////////
	_ "sdm_demo_todolist/sqlx/docs"
)

func Init(myRouter *gin.Engine) {

	// add swagger
	//
	// https://www.youtube.com/watch?v=AtaXj2hj074
	// -->
	// https://github.com/lemoncode21/golang-crud-gin-gorm

	//  https://hoohoo.top/blog/20210530112304-golang-tutorial-introduction-gin-html-template-and-how-integration-with-bootstrap/

	myRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// myRouter.StaticFile("/swagger/", "./swagger/index.html")
}

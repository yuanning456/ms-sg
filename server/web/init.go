package web

import (
	"ms-sg/db"

	"ms-sg/server/web/controller"
	"ms-sg/server/web/middleware"

	"github.com/gin-gonic/gin"
)

func Init(g *gin.Engine) {
	db.Init()
	initRouter(g)
}

func initRouter(g *gin.Engine) {
	g.Use(middleware.Cors())
	// g.Any("/account/register", controller.DefaultAccountController.Register)
	g.GET("/account/register", controller.DefaultAccountController.Register)
}
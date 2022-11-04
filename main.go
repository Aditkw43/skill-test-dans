package main

import (
	"skill-test-dans/config"
	"skill-test-dans/internal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	db := config.InitGorm()
	data := internal.GetData()
	controller := internal.NewController(db, data)

	public := r.Group("/")
	public.POST("register", controller.Register)
	public.POST("login", controller.Login)

	private := r.Group("/job")
	private.Use(internal.JwtAuthMiddleware())
	private.GET("/", controller.GetJobList)
	private.GET(":id", controller.GetJobDetail)

	logout := r.Group("/logout")
	logout.Use(internal.JwtAuthMiddleware())
	logout.POST("/", controller.Logout)

	r.Run(":8080")
}

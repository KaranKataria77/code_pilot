package main

import (
	"code_pilot/internal/config"
	"code_pilot/internal/middlewares"
	"code_pilot/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()

	r.SetTrustedProxies([]string{})

	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.CORSMiddleware())

	routes.SetUpRoutes(r)

	port := config.GetEnv("PORT", "8080")
	log.Println("Server running at PORT: ", port)
	r.Run(":" + port)
}

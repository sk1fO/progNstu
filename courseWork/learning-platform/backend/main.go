package main

import (
	"platform/db"
	"platform/handlers"
	"platform/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
		c.Set("db", database)
		c.Next()
	})

	// Публичные маршруты
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Приватные маршруты
	private := router.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/tasks", handlers.GetTasks)
		private.GET("/tasks/:id", handlers.GetTask)
		private.POST("/submit", handlers.SubmitSolution)
		private.GET("/submissions", handlers.GetSubmissions)
	}

	router.Run(":8080")
}

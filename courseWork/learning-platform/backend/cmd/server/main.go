package main

import (
	"learning-platform/backend/internal/config"
	"learning-platform/backend/internal/database"
	"learning-platform/backend/internal/handlers"
	"learning-platform/backend/internal/middleware"
	"learning-platform/backend/internal/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация Gin
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем базу данных
	database.InitDatabase(cfg)

	// Проверяем подключение к Docker
	_, err := services.NewDockerService()
	if err != nil {
		log.Printf("Warning: Docker service not available: %v", err)
	} else {
		log.Println("Docker service initialized")
	}

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/user", handlers.CurrentUser)
		protected.POST("/courses", handlers.CreateCourse)
		protected.GET("/courses", handlers.GetCourses)
		protected.POST("/lessons", handlers.CreateLesson)
		protected.GET("/courses/:courseId/lessons", handlers.GetLessons)
		protected.POST("/assignments", handlers.CreateAssignment)
		protected.GET("/lessons/:lessonId/assignments", handlers.GetAssignments)
		protected.POST("/assignments/:id/run", handlers.RunCode)
	}

	// Start server
	log.Println("Server starting on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

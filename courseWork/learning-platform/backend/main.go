package main

import (
	"log"
	"platform/db"
	"platform/handlers"
	"platform/middleware"
	"platform/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	// Инициализируем сервис заданий
	taskService := services.NewTaskService()
	if err := taskService.LoadTasksFromFile("data/tasks.json"); err != nil {
		log.Fatal("Ошибка загрузки заданий:", err)
	}

	// Инициализируем сервис валидации
	validationService := services.NewValidationService()

	// Синхронизируем задания с БД
	tasksForSync := taskService.GetTasksForSync()

	// Конвертируем в формат для БД
	var dbTasks []db.Task
	for _, task := range tasksForSync {
		dbTasks = append(dbTasks, db.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Difficulty:  task.Difficulty,
		})
	}

	if err := db.SyncTasks(database, dbTasks); err != nil {
		log.Fatal("Ошибка синхронизации заданий:", err)
	}

	// Передаем сервисы в обработчики
	handlers.SetTaskService(taskService)
	handlers.SetValidationService(validationService)

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

	log.Println("Сервер запущен на :8080")
	router.Run(":8080")
}

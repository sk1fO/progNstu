package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("your_secret_key") // В реальном приложении используйте переменные окружения

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=student teacher"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateCourse(c *gin.Context) {
	// Implementation for creating a course
	c.JSON(http.StatusOK, gin.H{"message": "Course creation endpoint"})
}

func GetCourses(c *gin.Context) {
	// Implementation for getting courses
	c.JSON(http.StatusOK, gin.H{"message": "Get courses endpoint"})
}

package handlers

import (
	"net/http"
	"strconv"

	"learning-platform/backend/internal/database"
	"learning-platform/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateLessonRequest struct {
	Title         string `json:"title" binding:"required"`
	Order         int    `json:"order" binding:"required"`
	TheoryContent string `json:"theory_content" binding:"required"`
	CourseID      uint   `json:"course_id" binding:"required"`
}

func CreateLesson(c *gin.Context) {
	username := c.GetString("username")

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers can create lessons"})
		return
	}

	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли курс и принадлежит ли он текущему пользователю
	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add lessons to your own courses"})
		return
	}

	lesson := models.Lesson{
		Title:         req.Title,
		Order:         req.Order,
		TheoryContent: req.TheoryContent,
		CourseID:      req.CourseID,
	}

	if err := database.DB.Create(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create lesson"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func GetLessons(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var lessons []models.Lesson
	if err := database.DB.Where("course_id = ?", courseID).Order("\"order\"").Find(&lessons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch lessons"})
		return
	}

	c.JSON(http.StatusOK, lessons)
}

type CreateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StarterCode string `json:"starter_code"`
	Language    string `json:"language" binding:"required"`
	LessonID    uint   `json:"lesson_id" binding:"required"`
}

func CreateAssignment(c *gin.Context) {
	username := c.GetString("username")

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers can create assignments"})
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли урок и принадлежит ли он курсу пользователя
	var lesson models.Lesson
	if err := database.DB.Preload("Course").First(&lesson, req.LessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	if lesson.Course.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add assignments to your own lessons"})
		return
	}

	assignment := models.Assignment{
		Title:       req.Title,
		Description: req.Description,
		StarterCode: req.StarterCode,
		Language:    req.Language,
		LessonID:    req.LessonID,
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create assignment"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func GetAssignments(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	var assignments []models.Assignment
	if err := database.DB.Where("lesson_id = ?", lessonID).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch assignments"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

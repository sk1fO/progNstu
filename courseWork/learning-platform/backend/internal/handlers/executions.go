package handlers

import (
	"net/http"
	"strconv"
	"time"

	"learning-platform/backend/internal/database"
	"learning-platform/backend/internal/models"
	"learning-platform/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ExecutionRequest struct {
	Code string `json:"code" binding:"required"`
}

func RunCode(c *gin.Context) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	username := c.GetString("username")
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var assignment models.Assignment
	if err := database.DB.Preload("TestCases").First(&assignment, assignmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	var req ExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем Docker сервис
	dockerService, err := services.NewDockerService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker service unavailable"})
		return
	}

	// Сохраняем решение в БД
	solution := models.Solution{
		Code:         req.Code,
		AssignmentID: uint(assignmentID),
		UserID:       user.ID,
		Status:       "sent",
	}
	if err := database.DB.Create(&solution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save solution"})
		return
	}

	// Запускаем код в Docker
	ctx := c.Request.Context() // Используем контекст запроса

	response, err := dockerService.RunCode(ctx, services.RunCodeRequest{
		Code:     req.Code,
		Language: assignment.Language,
		Timeout:  10 * time.Second,
	})

	if err != nil {
		solution.Status = "error"
		solution.TeacherComment = err.Error()
		database.DB.Save(&solution)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Execution failed",
			"details": err.Error(),
		})
		return
	}

	// Проверяем тесты
	testResults := make([]TestResult, 0)
	allPassed := true

	for _, test := range assignment.TestCases {
		// Для каждого теста запускаем проверку
		testResponse, err := dockerService.RunCode(ctx, services.RunCodeRequest{
			Code:     req.Code + "\n\n" + test.InputData,
			Language: assignment.Language,
			Timeout:  5 * time.Second,
		})

		if err != nil {
			testResults = append(testResults, TestResult{
				Passed:       false,
				ErrorMessage: err.Error(),
			})
			allPassed = false
			continue
		}

		isPassed := testResponse.Output == test.ExpectedOutput
		if !isPassed {
			allPassed = false
		}

		testResults = append(testResults, TestResult{
			Input:        test.InputData,
			Expected:     test.ExpectedOutput,
			Actual:       testResponse.Output,
			Passed:       isPassed,
			ErrorMessage: "",
		})
	}

	// Обновляем решение
	solution.Status = "tested"
	solution.PassedAutotests = allPassed
	database.DB.Save(&solution)

	c.JSON(http.StatusOK, gin.H{
		"solution_id":  solution.ID,
		"all_passed":   allPassed,
		"output":       response.Output,
		"test_results": testResults,
	})
}

type TestResult struct {
	Input        string `json:"input"`
	Expected     string `json:"expected"`
	Actual       string `json:"actual"`
	Passed       bool   `json:"passed"`
	ErrorMessage string `json:"error_message,omitempty"`
}

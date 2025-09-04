package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"platform/services"
	"strings"

	"github.com/gin-gonic/gin"
)

var validationService *services.ValidationService

func SetValidationService(service *services.ValidationService) {
	validationService = service
}

type SubmitRequest struct {
	TaskID int    `json:"task_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type SubmissionResponse struct {
	ID          int                   `json:"id"`
	TaskID      int                   `json:"task_id"`
	Code        string                `json:"code"`
	Status      string                `json:"status"`
	Output      string                `json:"output"`
	TestResults []services.TestResult `json:"test_results,omitempty"`
	CreatedAt   string                `json:"created_at"`
}

func SubmitSolution(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Получаем задание
	task, err := taskService.GetTaskByID(req.TaskID)
	if err != nil || task == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Задание не найдено"})
		return
	}

	// Проверяем, есть ли тесты у задания
	if task.Tests == nil || len(task.Tests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Для этого задания нет тестов"})
		return
	}

	// Конвертируем тесты в формат для валидации
	var validationTests []services.TestCase
	for _, test := range task.Tests {
		validationTests = append(validationTests, services.TestCase{
			Input:          test.Input,
			ExpectedOutput: test.ExpectedOutput,
			Description:    test.Description,
		})
	}

	// Запускаем валидацию тестов
	testResults, allPassed := validationService.ValidateSolution(req.Code, validationTests)

	status := "success"
	if !allPassed {
		status = "error"
	}

	// Формируем вывод для пользователя
	output := formatTestResults(testResults, allPassed)

	// Сохраняем результаты тестов в JSON
	testResultsJSON, err := json.Marshal(testResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения результатов тестов"})
		return
	}

	// Сохраняем решение в БД
	result, err := database.Exec(
		"INSERT INTO submissions (user_id, task_id, code, status, output, test_results) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.TaskID, req.Code, status, output, string(testResultsJSON),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения решения: " + err.Error()})
		return
	}

	submissionID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"message":       getStatusMessage(allPassed),
		"submission_id": submissionID,
		"status":        status,
		"output":        output,
		"test_results":  testResults,
		"passed":        allPassed,
	})
}

func GetSubmissions(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	rows, err := database.Query(`
        SELECT id, task_id, code, status, output, test_results, created_at 
        FROM submissions 
        WHERE user_id = ? 
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения решений: " + err.Error()})
		return
	}
	defer rows.Close()

	var submissions []SubmissionResponse
	for rows.Next() {
		var submission SubmissionResponse
		var testResultsJSON sql.NullString

		err := rows.Scan(
			&submission.ID,
			&submission.TaskID,
			&submission.Code,
			&submission.Status,
			&submission.Output,
			&testResultsJSON,
			&submission.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных: " + err.Error()})
			return
		}

		// Парсим результаты тестов если они есть
		if testResultsJSON.Valid && testResultsJSON.String != "" {
			var testResults []services.TestResult
			if err := json.Unmarshal([]byte(testResultsJSON.String), &testResults); err == nil {
				submission.TestResults = testResults
			}
		}

		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных: " + err.Error()})
		return
	}

	if submissions == nil {
		submissions = []SubmissionResponse{}
	}

	c.JSON(http.StatusOK, submissions)
}

func formatTestResults(results []services.TestResult, allPassed bool) string {
	if len(results) == 0 {
		return "Нет тестов для выполнения"
	}

	var output strings.Builder
	if allPassed {
		output.WriteString("✅ Все тесты пройдены успешно!\n\n")
	} else {
		output.WriteString("❌ Тесты не пройдены:\n\n")
	}

	for i, result := range results {
		output.WriteString(fmt.Sprintf("Тест %d: %s\n", i+1, result.Description))
		if result.Passed {
			output.WriteString("   ✅ Пройден\n")
		} else {
			output.WriteString("   ❌ Не пройден\n")
			if result.Input != "" {
				output.WriteString(fmt.Sprintf("   Вход: %s\n", result.Input))
			}
			output.WriteString(fmt.Sprintf("   Ожидалось: %s\n", result.Expected))
			output.WriteString(fmt.Sprintf("   Получено: %s\n", result.Actual))
		}
		output.WriteString("\n")
	}

	return output.String()
}

func getStatusMessage(passed bool) string {
	if passed {
		return "Решение верное! Все тесты пройдены."
	}
	return "Решение содержит ошибки."
}

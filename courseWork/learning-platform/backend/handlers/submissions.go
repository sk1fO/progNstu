package handlers

import (
	"database/sql"
	"net/http"
	"platform/models"
	"platform/utils"

	"github.com/gin-gonic/gin"
)

func SubmitSolution(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	var submission struct {
		TaskID int    `json:"task_id" binding:"required"`
		Code   string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Валидация: проверяем, что task_id существует
	var taskExists bool
	err := database.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", submission.TaskID).Scan(&taskExists)
	if err != nil || !taskExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Задание не найдено"})
		return
	}

	// Запускаем код в Docker
	output, err := utils.RunCppCode(submission.Code)
	status := "success"
	if err != nil {
		status = "error"
		output = err.Error()
	}

	// Сохраняем решение в БД
	result, err := database.Exec(
		"INSERT INTO submissions (user_id, task_id, code, status, output) VALUES (?, ?, ?, ?, ?)",
		userID, submission.TaskID, submission.Code, status, output,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения решения"})
		return
	}

	submissionID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"message":       "Решение проверено",
		"submission_id": submissionID,
		"status":        status,
		"output":        output,
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
        SELECT id, task_id, code, status, output 
        FROM submissions 
        WHERE user_id = ? 
        ORDER BY id DESC
    `, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения решений"})
		return
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var submission models.Submission
		if err := rows.Scan(
			&submission.ID,
			&submission.TaskID,
			&submission.Code,
			&submission.Status,
			&submission.Output,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

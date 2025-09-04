package handlers

import (
	"database/sql"
	"net/http"
	"platform/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	rows, err := database.Query("SELECT id, title, description, difficulty FROM tasks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заданий"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Difficulty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	taskID := c.Param("id")

	var task models.Task
	err := database.QueryRow("SELECT id, title, description, difficulty FROM tasks WHERE id = ?", taskID).Scan(
		&task.ID, &task.Title, &task.Description, &task.Difficulty)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задание не найдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

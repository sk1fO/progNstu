package handlers

import (
	"net/http"
	"platform/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

var taskService *services.TaskService

func SetTaskService(service *services.TaskService) {
	taskService = service
}

func GetTasks(c *gin.Context) {
	tasks := taskService.GetTasks()
	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задания"})
		return
	}

	task, err := taskService.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения задания"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задание не найдено"})
		return
	}

	c.JSON(http.StatusOK, task)
}

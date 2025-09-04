package services

import (
	"encoding/json"
	"fmt"
	"os"
	"platform/models"

	_ "github.com/mattn/go-sqlite3"
)

type TaskService struct {
	tasks []models.Task
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) GetTasks() []models.Task {
	return s.tasks
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	for _, task := range s.tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, nil
}

func (s *TaskService) GetTaskTests(taskID int) ([]models.TestCase, error) {
	task, err := s.GetTaskByID(taskID)
	if err != nil || task == nil {
		return nil, fmt.Errorf("задание не найдено")
	}
	return task.Tests, nil
}

// GetTasksForSync возвращает задачи в формате для синхронизации с БД
func (s *TaskService) GetTasksForSync() []struct {
	ID          int
	Title       string
	Description string
	Difficulty  string
} {
	var result []struct {
		ID          int
		Title       string
		Description string
		Difficulty  string
	}

	for _, task := range s.tasks {
		result = append(result, struct {
			ID          int
			Title       string
			Description string
			Difficulty  string
		}{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Difficulty:  task.Difficulty,
		})
	}

	return result
}

func (s *TaskService) LoadTasksFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var tasks []models.Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return err
	}

	// Отладочная информация
	for _, task := range tasks {
		fmt.Printf("Задание %d: %s, тестов: %d\n", task.ID, task.Title, len(task.Tests))
		for j, test := range task.Tests {
			fmt.Printf("  Тест %d: %s -> %s\n", j+1, test.Input, test.ExpectedOutput)
		}
	}

	s.tasks = tasks
	return nil
}

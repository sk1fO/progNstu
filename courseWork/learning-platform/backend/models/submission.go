package models

type Submission struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	TaskID int    `json:"task_id"`
	Code   string `json:"code"`
	Status string `json:"status"` // pending/success/error
	Output string `json:"output"` // вывод программы
}

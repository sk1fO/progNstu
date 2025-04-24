package data

import "fmt"

type Database interface {
	SaveAlert(buttonID string)
}

type MockDatabase struct{}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

func (db *MockDatabase) SaveAlert(buttonID string) {
	fmt.Printf("[MOCK DB] Событие от кнопки %s сохранено\n", buttonID)
}

package main

import (
	"monitoring/business"
	"monitoring/data"
	"monitoring/presentation"
)

func main() {
	// Инициализация слоёв
	db := data.NewMockDatabase()
	service := business.NewAlertService(db)
	ui := presentation.NewConsoleUI(service)

	// Запуск
	ui.Run()
}

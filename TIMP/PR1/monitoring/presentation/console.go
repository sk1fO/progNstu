package presentation

import (
	"bufio"
	"fmt"
	"os"

	"monitoring/business"
)

type ConsoleUI struct {
	service *business.AlertService
}

func NewConsoleUI(service *business.AlertService) *ConsoleUI {
	return &ConsoleUI{service: service}
}

func (ui *ConsoleUI) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Введите ID кнопки (или 'exit' для выхода): ")
		scanner.Scan()
		input := scanner.Text()
		if input == "exit" {
			break
		}
		ui.service.HandleAlert(input)
	}
}

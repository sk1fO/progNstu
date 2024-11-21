package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	file, _ := os.OpenFile("/home/sk1fo/vsCode/progNstu/Prac1Sem3/go/test/test.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Записываем любые буферизованные данные в подлежащий writer (стандартный вывод).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

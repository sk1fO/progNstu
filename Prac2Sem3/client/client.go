package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Print("Введите ip: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ip := scanner.Text()

	// Подключение к серверу
	conn, err := net.Dial("tcp", ip+":7432")
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу. Введите SQL-запросы или 'exit' для выхода.")

	reader := bufio.NewReader(os.Stdin)
	for {
		// Чтение запроса от пользователя
		fmt.Print("SQL> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Ошибка чтения ввода: %v", err)
		}

		// Удаление символа новой строки
		query = query[:len(query)-1]

		// Отправка запроса на сервер
		_, err = conn.Write([]byte(query + "\n"))
		if err != nil {
			log.Fatalf("Ошибка отправки запроса: %v", err)
		}

		// Чтение ответа от сервера
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatalf("Ошибка чтения ответа: %v", err)
		}

		// Удаление символа новой строки
		response = response[:len(response)-1]

		// Вывод ответа
		fmt.Println("Результат:")

		results := strings.Split(response, "]")
		for _, result := range results {
			if result == "" {
				break
			}

			result = strings.ReplaceAll(result, ",", " ")
			result = strings.ReplaceAll(result, `"`, "")
			fmt.Println(strings.ReplaceAll(result, "[", ""))
		}

		// Выход, если пользователь ввел "exit"
		if query == "exit" {
			fmt.Println("Выход из клиента.")
			return
		}
	}
}

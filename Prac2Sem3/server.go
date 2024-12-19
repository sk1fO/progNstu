package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Структура для хранения состояния сервера
type Server struct {
	db   *Database // База данных
	lock sync.Mutex
}

// Обработка подключения клиента
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Чтение запроса от клиента
		request, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка чтения запроса: %v", err)
			return
		}

		// Удаление символа новой строки
		request = request[:len(request)-1]

		// Обработка запроса
		s.lock.Lock()
		result, err := s.processQuery(request)
		s.lock.Unlock()

		if err != nil {
			result = fmt.Sprintf("Ошибка: %v", err)
		}

		// Отправка результата клиенту
		_, err = conn.Write([]byte(result + "\n"))
		if err != nil {
			log.Printf("Ошибка отправки ответа: %v", err)
			return
		}
	}
}

// Обработка SQL-запроса
func (s *Server) processQuery(query string) (string, error) {
	parts := strings.Split(query, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("пустой запрос")
	}

	switch parts[0] {
	case "exit":
		return "Выход из сервера", nil
	case "status":
		return "Сервер работает", nil
	case "INSERT":
		if len(parts) < 5 || parts[1] != "INTO" || parts[3] != "VALUES" {
			return "", fmt.Errorf("неверная команда INSERT. Использование: INSERT INTO <таблица> VALUES (<значения>)")
		}
		table := parts[2]
		values := strings.Trim(parts[4], "()")
		valuesList := strings.Split(values, ",")
		if err := insertData(s.db, table, valuesList); err != nil {
			return "", err
		}
		return "Данные успешно вставлены.", nil
	case "SELECT":
		if len(parts) < 3 || parts[2] != "FROM" {
			return "", fmt.Errorf("неверная команда SELECT. Использование: SELECT <столбцы> FROM <таблицы>")
		}
		result, err := selectData(s.db, query)
		if err != nil {
			return "", err
		}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(jsonResult), nil
	case "DELETE":
		if len(parts) < 4 || parts[1] != "FROM" {
			return "", fmt.Errorf("неверная команда DELETE. Использование: DELETE FROM <таблица> [WHERE <условия>]")
		}
		table := parts[2]
		conditions := ""
		if len(parts) > 4 && parts[3] == "WHERE" {
			conditions = strings.Join(parts[4:], " ")
		}
		if err := deleteData(s.db, table, conditions); err != nil {
			return "", err
		}
		return "Данные успешно удалены.", nil
	default:
		return "", fmt.Errorf("неизвестная команда: %s", parts[0])
	}
}

func main() {
	// Чтение конфигурации схемы
	schema, err := readSchema("schema.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать схему: %v", err)
	}

	// Создание базы данных
	db, err := createDatabase(schema)
	if err != nil {
		log.Fatalf("Не удалось создать базу данных: %v", err)
	}

	// Создание сервера
	server := &Server{
		db: db,
	}

	// Запуск сервера на порту 7432
	listener, err := net.Listen("tcp", ":7432")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
	defer listener.Close()

	log.Println("Сервер запущен на порту 7432")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при подключении клиента: %v", err)
			continue
		}

		// Обработка подключения в отдельной горутине
		go server.handleConnection(conn)
	}
}

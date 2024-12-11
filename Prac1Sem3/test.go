package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Определение структуры схемы
type Schema struct {
	Name       string              `json:"name"`         // Имя схемы
	TupleLimit int                 `json:"tuples_limit"` // Лимит строк в файле
	Structure  map[string][]string `json:"structure"`    // Структура таблиц
}

// Определение структуры таблицы
type Table struct {
	Name         string      // Имя таблицы
	Columns      []string    // Список столбцов
	PrimaryKey   string      // Первичный ключ
	SequenceFile string      // Файл последовательности для первичного ключа
	LockFile     string      // Файл блокировки
	Files        []string    // Список файлов CSV
	Lock         *sync.Mutex // Мьютекс для блокировки таблицы
}

// Определение структуры базы данных
type Database struct {
	Name       string            // Имя базы данных
	TupleLimit int               // Лимит строк в файле
	Tables     map[string]*Table // Список таблиц
}

// Чтение конфигурации схемы
func readSchema(filePath string) (*Schema, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

// Создание необходимых директорий и файлов на основе схемы
func createDatabase(schema *Schema) (*Database, error) {
	db := &Database{
		Name:       schema.Name,
		TupleLimit: schema.TupleLimit,
		Tables:     make(map[string]*Table),
	}

	baseDir := schema.Name
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	for tableName, columns := range schema.Structure {
		tableDir := filepath.Join(baseDir, tableName)
		if err := os.MkdirAll(tableDir, 0755); err != nil {
			return nil, err
		}

		primaryKey := fmt.Sprintf("%s_pk", tableName)
		sequenceFile := filepath.Join(tableDir, fmt.Sprintf("%s_sequence", tableName))
		lockFile := filepath.Join(tableDir, fmt.Sprintf("%s_lock", tableName))

		// Инициализация файла последовательности
		if _, err := os.Stat(sequenceFile); os.IsNotExist(err) {
			if err := os.WriteFile(sequenceFile, []byte("0"), 0644); err != nil {
				return nil, err
			}
		}

		// Инициализация файла блокировки
		if _, err := os.Stat(lockFile); os.IsNotExist(err) {
			if err := os.WriteFile(lockFile, []byte("false"), 0644); err != nil {
				return nil, err
			}
		}

		// Инициализация первого CSV файла
		csvFile := filepath.Join(tableDir, "1.csv")
		if _, err := os.Stat(csvFile); os.IsNotExist(err) {
			file, err := os.Create(csvFile)
			if err != nil {
				return nil, err
			}
			writer := csv.NewWriter(file)
			writer.Write(append([]string{primaryKey}, columns...))
			writer.Flush()
			file.Close()
		}

		db.Tables[tableName] = &Table{
			Name:         tableName,
			Columns:      append([]string{primaryKey}, columns...),
			PrimaryKey:   primaryKey,
			SequenceFile: sequenceFile,
			LockFile:     lockFile,
			Files:        []string{csvFile},
			Lock:         &sync.Mutex{},
		}
	}

	return db, nil
}

// Реализация операции SELECT
func selectData(db *Database, query string) ([][]string, error) {
	parts := strings.Split(query, " ")
	if len(parts) < 4 || parts[2] != "FROM" {
		return nil, fmt.Errorf("неверная команда SELECT. Использование: SELECT <столбцы> FROM <таблицы>")
	}

	selectCols := strings.Split(parts[1], ",")
	fromTables := strings.Split(parts[3], ",")

	var result [][]string

	for _, table := range fromTables {
		tbl, ok := db.Tables[table]
		if !ok {
			return nil, fmt.Errorf("таблица %s не существует", table)
		}

		for _, file := range tbl.Files {
			rows, err := readCSV(file)
			if err != nil {
				return nil, err
			}

			for _, row := range rows[1:] { // Пропускаем заголовок
				var selectedRow []string
				for _, col := range selectCols {
					colParts := strings.Split(col, ".")
					if colParts[0] == table {
						idx := -1
						for i, c := range tbl.Columns {
							if c == colParts[1] {
								idx = i
								break
							}
						}
						if idx != -1 {
							selectedRow = append(selectedRow, row[idx])
						}
					}
				}
				result = append(result, selectedRow)
			}
		}
	}

	return result, nil
}

// Реализация операции WHERE
func whereClause(rows [][]string, conditions string) [][]string {
	var header []string = rows[0]
	var result [][]string
	for _, row := range rows {
		if evaluateCondition(header, row, conditions) {
			result = append(result, row)
		}
	}
	return result
}

// Оценка условия WHERE
func evaluateCondition(header, row []string, conditions string) bool {
	// Простая оценка равенства
	parts := strings.Split(conditions, "=")
	col := parts[0]
	value := parts[1]

	colParts := strings.Split(col, ".")
	idx := -1
	for i, c := range header { // из заголовка!!!
		if c == colParts[1] {
			idx = i
			break
		}
	}

	if idx != -1 {
		return row[idx] == value
	}

	return false
}

// Реализация операции INSERT
func insertData(db *Database, table string, values []string) error {
	tbl, ok := db.Tables[table]
	if !ok {
		return fmt.Errorf("таблица %s не существует", table)
	}

	tbl.Lock.Lock()
	defer tbl.Lock.Unlock()

	// Получение следующего первичного ключа
	seq, err := getNextSequence(tbl.SequenceFile)
	if err != nil {
		return err
	}

	// Добавление первичного ключа к значениям
	values = append([]string{strconv.Itoa(seq)}, values...)

	// Определение файла, в который нужно записать данные
	lastFile := tbl.Files[len(tbl.Files)-1]
	rows, err := readCSV(lastFile)
	if err != nil {
		return err
	}

	if len(rows) >= db.TupleLimit {
		// Создание нового файла
		newFile := filepath.Join(filepath.Dir(lastFile), fmt.Sprintf("%d.csv", len(tbl.Files)+1))
		file, err := os.Create(newFile)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		writer.Write(tbl.Columns)
		writer.Write(values)
		writer.Flush()
		file.Close()
		tbl.Files = append(tbl.Files, newFile)
	} else {
		// Запись в последний файл
		file, err := os.OpenFile(lastFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		writer.Write(values)
		writer.Flush()
		file.Close()
	}

	return nil
}

// Получение следующего номера последовательности
func getNextSequence(file string) (int, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}
	seq, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	seq++
	if err := os.WriteFile(file, []byte(strconv.Itoa(seq)), 0644); err != nil {
		return 0, err
	}
	return seq, nil
}

// Реализация операции DELETE
func deleteData(db *Database, table string, conditions string) error {
	tbl, ok := db.Tables[table]
	if !ok {
		return fmt.Errorf("таблица %s не существует", table)
	}

	tbl.Lock.Lock()
	defer tbl.Lock.Unlock()

	for i, file := range tbl.Files {
		rows, err := readCSV(file)
		if err != nil {
			return err
		}

		header := rows[0]
		newRows := [][]string{header} // Оставляем заголовок

		for _, row := range rows[1:] {
			if !evaluateCondition(header, row, conditions) {
				newRows = append(newRows, row)
			}
		}

		// Если файл становится пустым после удаления, удаляем его
		if len(newRows) == 1 {
			if err := os.Remove(file); err != nil {
				return err
			}
			// Удаляем файл из списка
			tbl.Files = append(tbl.Files[:i], tbl.Files[i+1:]...)
		} else {
			// Запись новых строк обратно в файл
			file, err := os.Create(file)
			if err != nil {
				return err
			}
			writer := csv.NewWriter(file)
			for _, row := range newRows {
				if err := writer.Write(row); err != nil {
					file.Close()
					return err
				}
			}
			writer.Flush()
			file.Close()
		}
	}

	return nil
}

// Чтение CSV файла
func readCSV(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// Главная функция для запуска базы данных
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

	// Консольное меню для ввода команд SQL
	for {
		fmt.Print("\nВведите команду SQL (или 'exit' для выхода): ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()

		if command == "exit" {
			fmt.Println("Выход...")
			return
		}

		parts := strings.Split(command, " ")

		switch parts[0] {
		case "INSERT":
			if len(parts) < 5 || parts[1] != "INTO" || parts[3] != "VALUES" {
				fmt.Println("Неверная команда INSERT. Использование: INSERT INTO <таблица> VALUES (<значения>)")
				continue
			}
			table := parts[2]
			values := strings.Trim(parts[4], "()")
			valuesList := strings.Split(values, ",")
			if err := insertData(db, table, valuesList); err != nil {
				fmt.Println("Ошибка при вставке данных:", err)
			} else {
				fmt.Println("Данные успешно вставлены.")
			}
		case "SELECT":
			if len(parts) < 3 || parts[2] != "FROM" {
				fmt.Println("Неверная команда SELECT. Использование: SELECT <столбцы> FROM <таблицы>")
				continue
			}
			query := strings.Join(parts, " ")
			result, err := selectData(db, query)
			if err != nil {
				fmt.Println("Ошибка при выборке данных:", err)
			} else {
				for _, row := range result {
					fmt.Println(row)
				}
			}
		case "DELETE":
			if len(parts) < 4 || parts[1] != "FROM" {
				fmt.Println("Неверная команда DELETE. Использование: DELETE FROM <таблица> [WHERE <условия>]")
				continue
			}
			table := parts[2]
			conditions := ""
			if len(parts) > 4 && parts[3] == "WHERE" {
				conditions = strings.Join(parts[4:], " ")
			}
			if err := deleteData(db, table, conditions); err != nil {
				fmt.Println("Ошибка при удалении данных:", err)
			} else {
				fmt.Println("Данные успешно удалены.")
			}
		default:
			fmt.Println("Неизвестная команда. Поддерживаемые команды: INSERT, SELECT, DELETE, exit")
		}
	}
}

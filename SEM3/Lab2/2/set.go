package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Структура множества, основанного на хеш-таблице
type Set struct {
	hashTable *HashTable
}

// Создает и возвращает указатель на новое множество
func NewSet() *Set {
	return &Set{
		hashTable: NewHashTable(),
	}
}

// Добавляет элемент в множество
func (s *Set) Add(key string) {
	s.hashTable.Set(key, true) // Используем ключ и значение true
}

// Удаляет элемент из множества
func (s *Set) Remove(key string) {
	s.hashTable.Delete(key)
}

// Проверяет наличие элемента в множестве
func (s *Set) Contains(key string) bool {
	_, exists := s.hashTable.Get(key)
	return exists
}

// Возвращает все элементы множества в виде среза строк
func (s *Set) Read() []string {
	elements := s.hashTable.Read()
	result := make([]string, 0, len(elements))
	for key := range elements {
		result = append(result, key)
	}
	return result
}

// Основная функция для обработки входных данных и запросов
func main() {
	// Парсим аргументы командной строки
	fileFlag := flag.String("file", "", "Путь до файла с данными")
	queryPtr := flag.String("query", "", "Запрос к файлу с данными")
	flag.Parse()

	if *fileFlag == "" || *queryPtr == "" {
		log.Fatalln("Необходимо указать путь до файла и запрос.")
	}

	// Читаем данные из файла
	data, err := os.ReadFile(*fileFlag)
	if err != nil {
		log.Fatalln("Ошибка при чтении файла: ", err)
	}

	// Создаем множество
	set := NewSet()

	// Парсим данные из файла и добавляем их в множество
	var elements []string
	err = json.Unmarshal(data, &elements)
	if err != nil {
		fmt.Printf("Ошибка при парсинге данных из файла: %v\n", err)
		os.Exit(1)
	}

	for _, element := range elements {
		set.Add(element)
	}

	parts := strings.Split(*queryPtr, " ") // Разделяем команду на части
	command := parts[0]                    // Команда
	args := parts[1:]                      // Аргументы команды

	// Обрабатываем запрос
	switch command {
	case "PRINT":
		handlePrint(set)
	case "SET_AT":
		handleSetAt(args, set)
	case "SETADD":
		handleSetAdd(args, set)
	case "SETDEL":
		handleSetDel(args, set)
	default:
		log.Fatal("Неизвестный запрос.", command)
	}
	saveSetToFile(set, *fileFlag)
}

// Функция для сохранения множества в файл
func saveSetToFile(set *Set, filePath string) {
	elements := set.Read()
	data, err := json.Marshal(elements)
	if err != nil {
		log.Fatalln("Ошибка при сериализации данных: ", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Fatalln("Ошибка при записи данных в файл: ", err)
	}

}

func handlePrint(set *Set) {
	elements := set.Read()
	fmt.Println(elements)

}

func handleSetAt(args []string, s *Set) {
	if len(args) != 1 {
		log.Fatal("неверный запрос", args)
	}

	if s.Contains(args[0]) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func handleSetAdd(args []string, s *Set) {
	if len(args) != 1 {
		log.Fatal("неверный запрос", args)
	}

	s.Add(args[0])
}

func handleSetDel(args []string, s *Set) {
	if len(args) != 1 {
		log.Fatal("неверный запрос", args)
	}

	s.Remove(args[0])
}

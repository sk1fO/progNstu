package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const ABSOLUTE_PATH string = "/home/sk1fo/vsCode/progNstu/Prac1Sem3/go/test/"

type dbNode struct {
	value []string
	next  *dbNode
}

type db struct {
	len  int
	head *dbNode
}

// инициализация базы данных
func dbNew() *db {
	return &db{}
}

// деструктор
func (db *db) delete() { // если на объект нет указателй
	db.head = nil // то сборщик мусора удалит объект
	db.len = 0    // поэтому можно просто хед = нил
}

// добавление в конец
func (db *db) Push(_value []string) {
	node := &dbNode{ // создаем новый узел
		value: _value, // присваиваем значение
	} // указатель на некст по умолчанию нил

	if db.head == nil { // если узел единственный
		db.head = node // то он и есть хед
	} else { // иначе ищем последний узел
		current := db.head
		for current.next != nil {
			current = current.next
		}
		current.next = node // присваиваем ему значение
	}
	db.len++
}

// удаление последнего элемента
func (db *db) Pop() error {
	if db.head == nil { // обрабатываем ошибку
		return fmt.Errorf("database is empty")
	}

	var prev *dbNode   // сохраняем указатель на прошлый узел
	current := db.head //

	for current.next != nil { // в цикле ищем последний
		prev = current
		current = current.next
	}

	if prev != nil { // переподвязываем узлы
		prev.next = nil
	} else {
		db.head = nil
	}
	db.len-- // уменьшаем длину
	return nil
}

// функция распечатки
func (db *db) Print() { // нужна для дебага
	if db.head == nil { // проверка на пустоту
		fmt.Println("database is empty")
		return
	}

	current := db.head
	for current != nil { // в цикле проходим все значения
		fmt.Println(current.value) // печатаем их
		current = current.next
	}
}

// функция записи стркутуры в csv файл
func (db *db) writeCSV(path string) error {
	if db.head == nil {
		return fmt.Errorf("database is empty")
	}

	var toWrite [][]string
	current := db.head
	for current != nil {
		toWrite = append(toWrite, current.value)
		current = current.next
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(toWrite)
	return nil
}

// функция чтения структуры из файла
func (db *db) readCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	recs, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, rec := range recs {
		db.Push(rec)
	}
	return nil
}

// Config представляет конфигурацию из JSON-файла
type Config struct {
	Name        string              `json:"name"`
	TuplesLimit int                 `json:"tuples_limit"`
	Structure   map[string][]string `json:"structure"`
}

var config Config //глобальная переменная структуры субд

func ReadJson() error {
	// Открываем JSON-файл
	jsonFile, err := os.Open(ABSOLUTE_PATH + "scheme.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Читаем содержимое файла
	byteValue, _ := io.ReadAll(jsonFile)

	// Десериализуем JSON в структуру Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}

	// Выводим структуру данных
	fmt.Printf("Имя: %s\n", config.Name)
	fmt.Printf("Лимит кортежей: %d\n", config.TuplesLimit)

	fmt.Println("Структура таблиц:")
	for tableName, columns := range config.Structure {
		fmt.Printf("Таблица: %s\n", tableName)
		fmt.Print("\t")
		for _, column := range columns {
			fmt.Printf("%s ", column)
		}
		fmt.Println()
	}
	return nil
}

// функция создания директорий, принимает данные из config
func CreateDir() {
	os.Mkdir(ABSOLUTE_PATH+config.Name, 0777) //0777 -rwx
	for tableName := range config.Structure {
		os.Mkdir(ABSOLUTE_PATH+config.Name+"/"+tableName, 0777)
	}
}

func main() {
	val := []string{"1", "lox", "lol"}
	// val2 := []string{"2", "dab", "aga"}
	// val3 := []string{"3", "what&", "killer"}
	database := dbNew()
	database2 := dbNew()
	//
	// database.Push(val2)
	// database.Push(val3)

	database.Print()

	err := database.readCSV("/home/sk1fo/vsCode/progNstu/Prac1Sem3/go/test/test.csv")
	fmt.Println(err)

	database.Print()

	//database.delete()
	database.Push(val)
	database2.Print()
	database.Print()
}

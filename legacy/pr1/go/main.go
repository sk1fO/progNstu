package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const ABSOLUTE_PATH string = "/home/sk1fo/vsCode/progNstu/Prac1Sem3/go/"

// Config представляет конфигурацию из JSON-файла
type Config struct {
	Name        string              `json:"name"`
	TuplesLimit int                 `json:"tuples_limit"`
	Structure   map[string][]string `json:"structure"`
}

var config Config //глобальная переменная структуры субд

func ReadJson() {
	// Открываем JSON-файл
	jsonFile, err := os.Open(ABSOLUTE_PATH + "scheme.json")
	if err != nil {
		log.Fatalf("Не удалось открыть JSON: %v", err)
	}
	defer jsonFile.Close()

	// Читаем содержимое файла
	byteValue, _ := io.ReadAll(jsonFile)

	// Десериализуем JSON в структуру Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatalf("Не удалось десериализовать JSON: %v", err)
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
}

// функция создания директорий, принимает данные из config
func CreateDir() {
	os.Mkdir(ABSOLUTE_PATH+config.Name, 0777) //0777 -rwx
	for tableName := range config.Structure {
		os.Mkdir(ABSOLUTE_PATH+config.Name+"/"+tableName, 0777)
	}
}

func CreatePath(tableName string) string {
	return ABSOLUTE_PATH + config.Name + "/" + tableName
}

// функция поиска последнего .csv, принимает: название таблицы
func pathToMax(_table string) string {
	path := CreatePath(_table)
	// считываем список файлов в директории
	lst, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// обработка случая, когда нет файлов .csv
	if len(lst) == 0 {
		// создаем первый файл
		os.Create(path + "/1.csv")
		// открываем его на запись
		file, _ := os.OpenFile(path+"/1.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		writer := csv.NewWriter(file)

		// читаем имя таблицы
		table := strings.Split(path, "/")
		table_name := table[len(table)-1]

		// пишем хедер таблицы из JSON
		to_write := config.Structure[table_name]
		writer.Write(to_write)
		writer.Flush()
		file.Close()

		// возвращаем путь к файлу
		return path + "/1.csv"
	}
	//возвращает абсолютный путь к последнему файлу
	return path + "/" + lst[len(lst)-1].Name()
}

// реализация вставки, принимает: название таблицы, слайс значений
func INSERT_INTO(table string, value []string) {
	//открываем файл на чтение
	path := CreatePath(table)
	file, err := os.Open(pathToMax(table))
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file) //парсинг pk...

	values, _ := reader.ReadAll()
	splitPath := strings.Split(pathToMax(table), "/") //сплит пути к файлу
	maxNumStr := splitPath[len(splitPath)-1]
	maxNumInt, _ := strconv.Atoi(maxNumStr[:len(maxNumStr)-4]) //приводим к int
	table_name := splitPath[len(splitPath)-2]                  //имя таблицы
	pk_max, _ := strconv.Atoi(values[len(values)-1][0])        //получаем пк

	pk_max++ //прибавляем пк

	//проверка на максимальное кол-во записей в файле
	if pk_max == config.TuplesLimit*maxNumInt {
		maxNumInt++
		file.Close()
		//создаем файл со следующим по порядку номером
		file, err = os.Create(path + "/" + strconv.Itoa(maxNumInt) + ".csv")
		//открываем этот файл в режиме записи в конец
		file, err = os.OpenFile(path+"/"+strconv.Itoa(maxNumInt)+".csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}

		//записываем хедер таблицы
		writer := csv.NewWriter(file)
		to_write := config.Structure[table_name]
		writer.Write(to_write)
		writer.Flush()
		file.Close()

	}

	//открываем файл на запись в конец
	file, err = os.OpenFile(pathToMax(table), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	writer := csv.NewWriter(file)

	//формируем срез, куда запишем данные
	to_write := []string{strconv.Itoa(pk_max)}
	to_write = append(to_write, value...)

	//записываем всё в файл
	writer.Write(to_write)
	writer.Flush()
	file.Close()
}

func cross_join() {

}

func SELECT(table_column []string, tables []string) {

	temp := dbNew()
	//selected := dbNew()

	if len(tables) > 1 {
		cross_join()
	} else {
		path := CreatePath(tables[0])
		// считываем список файлов в директории
		lst, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range lst {
			temp.readCSV(path + "/" + file.Name())
		}
		temp.Print()
	}
}

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

func init() {
	ReadJson()
	CreateDir()
}
func main() {
	//value := []string{"51123", "Пенталгин", "Парацетамол", "23"}
	//INSERT_INTO("лекарства", value)
	sel := []string{"лекарства.вещество"}
	tab := []string{"лекарства"}
	SELECT(sel, tab)
}

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

// функция поиска последнего .csv
func pathToMax(path string) string {

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

// реализация вставки, принимает: путь к таблице, слайс значений
func INSERT_INTO(path string, value []string) {
	//открываем файл на чтение
	file, err := os.Open(pathToMax(path))
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file) //парсинг pk...

	values, _ := reader.ReadAll()
	splitPath := strings.Split(pathToMax(path), "/") //сплит пути к файлу
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
	file, err = os.OpenFile(pathToMax(path), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	writer := csv.NewWriter(file)

	//формируем срез, куда запишем данные
	to_write := []string{strconv.Itoa(pk_max)}
	to_write = append(to_write, value...)

	//записываем всё в файл
	writer.Write(to_write)
	writer.Flush()
	file.Close()
}

func init() {
	ReadJson()
	CreateDir()
}
func main() {
	value := []string{"89538834111", "Игорь", "15%"}
	INSERT_INTO(CreatePath("лояльность"), value)
}

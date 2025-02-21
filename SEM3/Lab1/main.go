package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var ( // map для записи в json
	arrays      = make(map[string]*Array) // объявляем все атд
	lists       = make(map[string]*List)
	singlelists = make(map[string]*SingleLinkedList)
	queues      = make(map[string]*Queue)
	stacks      = make(map[string]*Stack)
	hashtables  = make(map[string]*HashTable)
	cbtrees     = make(map[string]*CBTree)
	names       = make(map[string]interface{}) // Храним все имена с типом, которому всё присвоено
)

func main() {
	filePtr := flag.String("file", "", "file to read and write data") // Указатель на файл для чтения и записи данных
	queryPtr := flag.String("query", "", "query to execute")          // Указатель на команду для выполнения
	flag.Parse()                                                      // Парсим аргументы командной строки

	if *filePtr == "" || *queryPtr == "" {
		log.Fatal("Both --file and --query must be specified") // Проверяем, что указаны оба параметра
	}

	loadFromFile(*filePtr) // Загружаем данные из файла

	parts := strings.Split(*queryPtr, " ") // Разделяем команду на части
	command := parts[0]                    // Команда
	args := parts[1:]                      // Аргументы команды

	switch command {
	case "APUSH":
		handleArraySet(args) // Добавление элемента в массив
	case "APUSHINDEX":
		handleArraySetIndex(args)
	case "AGET":
		handleArrayGet(args) // Получение элемента из массива
	case "ADEL":
		handleArrayRemove(args) // Удаление элемента из массива
	case "AREPLACE":
		handleArrayReplace(args) // Замена элемента в массиве
	case "ALEN":
		handleArrayLength(args) // Получение длины массива

	case "LPUSHHEAD":
		handleListAddHead(args) // Добавление элемента в голову списка
	case "LPUSHTAIL":
		handleListAddTail(args) // Добавление элемента в хвост списка
	case "LDELHEAD":
		handleListRemoveHead(args) // Удаление элемента из головы списка
	case "LDELTAIL":
		handleListRemoveTail(args) // Удаление элемента из хвоста списка
	case "LDELVALUE":
		handleListRemoveValue(args) // Удаление элемента по значению
	case "LFINDVALUE":
		handleListFindValue(args) // Наличие элемента по значению

	case "SLPUSHHEAD":
		handleSListAddHead(args) // Добавление элемента в голову списка
	case "SLPUSHTAIL":
		handleSListAddTail(args) // Добавление элемента в хвост списка
	case "SLDELHEAD":
		handleSListRemoveHead(args) // Удаление элемента из головы списка
	case "SLDELTAIL":
		handleSListRemoveTail(args) // Удаление элемента из хвоста списка
	case "SLDELVALUE":
		handleSListRemoveValue(args) // Удаление элемента по значению
	case "SLFINDVALUE":
		handleSListFindValue(args) // Наличие элемента по значению

	case "QPUSH":
		handleQueuePush(args) // Добавление элемента в очередь
	case "QPOP":
		handleQueuePop(args) // Удаление элемента из очереди

	case "SPUSH":
		handleStackPush(args) // Добавление элемента в стек
	case "SPOP":
		handleStackPop(args) // Удаление элемента из стека

	case "HSET":
		handleHashSet(args) // Добавление элемента в хеш-таблицу
	case "HGET":
		handleHashGet(args) // Получение элемента из хеш-таблицы
	case "HDEL":
		handleHashDelete(args) // Удаление элемента из хеш-таблицы

	case "TINSERT":
		handleCBTAdd(args) // Добавление элемента в полное бинарное дерево
	case "TFIND":
		handleCBTFind(args) // Поиск элемента в полном бинарном дереве
	case "TISCOMPLETE":
		handleCBTIsComplete(args) // Проверка, является ли дерево полным

	case "PRINT":
		handlePrint(args) // Вывод структуры данных на экран

	default:
		log.Fatalf("Unknown command: %s", command) // Обработка неизвестной команды
	}

	saveToFile(*filePtr) // Сохраняем данные в файл
}

func handleArraySet(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: APUSH <array_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if arr, ok := arrays[name]; ok {
		arr.AddToEnd(value) // Добавляем элемент в массив по индексу
	} else {
		arr := NewArray()
		names[name] = "array"
		arr.AddToEnd(value) // Создаем новый массив и добавляем элемент
		arrays[name] = arr
	}
}

func handleArraySetIndex(args []string) {
	if len(args) < 3 {
		log.Fatal("Usage: APUSH <array_name> <index> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Usage: APUSH <array_name> <index> <value>\nIndex invalid: ", args[1]) // Проверка на корректность индекса
	}

	value := args[2]

	if arr, ok := arrays[name]; ok {
		arr.AddAtIndex(index, value) // Добавляем элемент в массив по индексу
	} else {
		arr := NewArray()
		names[name] = "array"
		arr.AddAtIndex(index, value) // Создаем новый массив и добавляем элемент
		arrays[name] = arr
	}
}

func handleArrayGet(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: AGET <array_name> <index>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Usage: AGET <array_name> <index>\nInvalid index: ", args[1]) // Проверка на корректность индекса
	}

	if arr, ok := arrays[name]; ok {
		value, err := arr.Get(index) // Получаем элемент из массива
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		fmt.Println(value)
	} else {
		log.Fatalf("Array %s not found", name) // Обработка случая, когда массив не найден
	}
}

func handleArrayRemove(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: ADEL <array_name> <index>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Usage: ADEL <array_name> <index>\nInvalid index: ", args[1]) // Проверка на корректность индекса
	}

	if arr, ok := arrays[name]; ok {
		err := arr.RemoveAtIndex(index) // Удаляем элемент из массива по индексу
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("Array %s not found", name) // Обработка случая, когда массив не найден
	}
}

func handleArrayReplace(args []string) {
	if len(args) < 3 {
		log.Fatal("Usage: AREPLACE <array_name> <index> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Usage: AREPLACE <array_name> <index> <value>\nInvalid index: ", args[1]) // Проверка на корректность индекса
	}

	value := args[2]

	if arr, ok := arrays[name]; ok {
		arr.ReplaceAtIndex(index, value) // Заменяем элемент в массиве по индексу
	} else {
		log.Fatalf("Array %s not found", name) // Обработка случая, когда массив не найден
	}
}

func handleArrayLength(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: ALEN <array_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if arr, ok := arrays[name]; ok {
		fmt.Println(arr.Length()) // Выводим длину массива
	} else {
		log.Fatalf("Array %s not found", name) // Обработка случая, когда массив не найден
	}
}

func handleSListAddHead(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: SLPUSHHEAD <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := singlelists[name]; ok {

		lst.AddToHead(value) // Добавляем элемент в голову списка
	} else {
		lst := NewSingleList()
		names[name] = "sl"
		lst.AddToHead(value) // Создаем новый список и добавляем элемент
		singlelists[name] = lst
	}
}

func handleSListAddTail(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: SLPUSHTAIL <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := singlelists[name]; ok {
		lst.AddToTail(value) // Добавляем элемент в хвост списка
	} else {
		lst := NewSingleList()
		names[name] = "sl"
		lst.AddToTail(value) // Создаем новый список и добавляем элемент
		singlelists[name] = lst
	}
}

func handleSListRemoveHead(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: SLDELHEAD <list_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if lst, ok := singlelists[name]; ok {
		value, err := lst.RemoveFromHead() // Удаляем элемент из головы списка
		if err != nil {                    // обработка ошибки
			log.Fatalln(err)
		}
		fmt.Println(value)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleSListRemoveTail(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: SLDELTAIL <list_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if lst, ok := singlelists[name]; ok {
		value, err := lst.RemoveFromTail() // Удаляем элемент из хвоста списка
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(value)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleSListRemoveValue(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: SLDELVALUE <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := singlelists[name]; ok {
		err := lst.RemoveByValue(value) // Удаляем элемент по значению
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleSListFindValue(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: SLFINDVALUE <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := singlelists[name]; ok {
		found := lst.FindByValue(value) // Ищем элемент по значению
		fmt.Println(found)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleListAddHead(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: LPUSHHEAD <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := lists[name]; ok {
		lst.AddToHead(value) // Добавляем элемент в голову списка
	} else {
		lst := NewList()
		names[name] = "list"
		lst.AddToHead(value) // Создаем новый список и добавляем элемент
		lists[name] = lst
	}
}

func handleListAddTail(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: LPUSHTAIL <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := lists[name]; ok {
		lst.AddToTail(value) // Добавляем элемент в хвост списка
	} else {
		lst := NewList()
		names[name] = "list"
		lst.AddToTail(value) // Создаем новый список и добавляем элемент
		lists[name] = lst
	}
}

func handleListRemoveHead(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: LDELHEAD <list_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if lst, ok := lists[name]; ok {
		value, _ := lst.RemoveFromHead() // Удаляем элемент из головы списка
		fmt.Println(value)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleListRemoveTail(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: LDELTAIL <list_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if lst, ok := lists[name]; ok {
		value, _ := lst.RemoveFromTail() // Удаляем элемент из хвоста списка
		fmt.Println(value)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleListRemoveValue(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: LDELVALUE <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := lists[name]; ok {
		lst.RemoveByValue(value) // Удаляем элемент по значению
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleListFindValue(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: LFINDVALUE <list_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if lst, ok := lists[name]; ok {
		found := lst.FindByValue(value) // Ищем элемент по значению
		fmt.Println(found)
	} else {
		log.Fatalf("List %s not found", name) // Обработка случая, когда список не найден
	}
}

func handleQueuePush(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: QPUSH <queue_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if q, ok := queues[name]; ok {
		q.Push(value) // Добавляем элемент в очередь
	} else {
		q := NewQueue()
		names[name] = "queue"
		q.Push(value) // Создаем новую очередь и добавляем элемент
		queues[name] = q
	}
}

func handleQueuePop(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: QPOP <queue_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if q, ok := queues[name]; ok {
		value, _ := q.Pop() // Удаляем элемент из очереди
		fmt.Println(value)
	} else {
		log.Fatalf("Queue %s not found", name) // Обработка случая, когда очередь не найдена
	}
}

func handleStackPush(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: SPUSH <stack_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if s, ok := stacks[name]; ok {
		s.Push(value) // Добавляем элемент в стек
	} else {
		s := NewStack()
		names[name] = "stack"
		s.Push(value) // Создаем новый стек и добавляем элемент
		stacks[name] = s
	}
}

func handleStackPop(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: SPOP <stack_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if s, ok := stacks[name]; ok {
		value, _ := s.Pop() // Удаляем элемент из стека
		fmt.Println(value)
	} else {
		log.Fatalf("Stack %s not found", name) // Обработка случая, когда стек не найден
	}
}

func handleHashSet(args []string) {
	if len(args) < 3 {
		log.Fatal("Usage: HSET <hash_name> <key> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	key := args[1]
	value := args[2]

	if h, ok := hashtables[name]; ok {
		h.Set(key, value) // Добавляем элемент в хеш-таблицу
	} else {
		h := NewHashTable()
		names[name] = "ht"
		h.Set(key, value) // Создаем новую хеш-таблицу и добавляем элемент
		hashtables[name] = h
	}
}

func handleHashGet(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: HGET <hash_name> <key>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	key := args[1]

	if h, ok := hashtables[name]; ok {
		value, found := h.Get(key) // Получаем элемент из хеш-таблицы
		if found {
			fmt.Println(value)
		} else {
			fmt.Println("Key not found")
		}
	} else {
		log.Fatalf("Hash table %s not found", name) // Обработка случая, когда хеш-таблица не найдена
	}
}

func handleHashDelete(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: HDEL <hash_name> <key>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	key := args[1]

	if h, ok := hashtables[name]; ok {
		h.Delete(key) // Удаляем элемент из хеш-таблицы
	} else {
		log.Fatalf("Hash table %s not found", name) // Обработка случая, когда хеш-таблица не найдена
	}
}

func handleCBTAdd(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: TINSERT <cbt_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if t, ok := cbtrees[name]; ok {
		t.Add(value) // Добавляем элемент в полное бинарное дерево
	} else {
		t := NewCBTree()
		names[name] = "cbt"
		t.Add(value) // Создаем новое полное бинарное дерево и добавляем элемент
		cbtrees[name] = t
	}
}

func handleCBTFind(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: TFIND <cbt_name> <value>") // Проверка на правильное количество аргументов
	}
	name := args[0]
	value := args[1]

	if t, ok := cbtrees[name]; ok {
		found := t.Find(value) // Ищем элемент в полном бинарном дереве
		fmt.Println(found)
	} else {
		log.Fatalf("CBTree %s not found", name) // Обработка случая, когда дерево не найдено
	}
}

func handleCBTIsComplete(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: TISCOMPLETE <cbt_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if t, ok := cbtrees[name]; ok {
		complete := t.IsComplete() // Проверяем, является ли дерево полным
		fmt.Println(complete)
	} else {
		log.Fatalf("CBTree %s not found", name) // Обработка случая, когда дерево не найдено
	}
}

func handlePrint(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: PRINT <structure_name>") // Проверка на правильное количество аргументов
	}
	name := args[0]

	if arr, ok := arrays[name]; ok {
		fmt.Println(arr.Read()) // Выводим массив
	} else if sl, ok := singlelists[name]; ok {
		fmt.Println(sl.Read()) // выводим список
	} else if lst, ok := lists[name]; ok {
		fmt.Println(lst.Read()) // Выводим список
	} else if q, ok := queues[name]; ok {
		fmt.Println(q.Read()) // Выводим очередь
	} else if s, ok := stacks[name]; ok {
		fmt.Println(s.Read()) // Выводим стек
	} else if h, ok := hashtables[name]; ok {
		fmt.Println(h.Read()) // Выводим хеш-таблицу
	} else if t, ok := cbtrees[name]; ok {
		fmt.Println(t.Read()) // Выводим полное бинарное дерево
	} else {
		log.Fatalf("Structure %s not found", name) // Обработка случая, когда структура данных не найдена
	}
}

func loadFromFile(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err) // Обработка ошибки чтения файла
		return
	}

	var savedData map[string]interface{}
	if err := json.Unmarshal(data, &savedData); err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err) // Обработка ошибки десериализации JSON
		return
	}

	names, _ = savedData["names"].(map[string]interface{})
	delete(savedData, "names")

	for k, v := range savedData {

		name := names[k].(string)
		switch name {
		case "array":
			arr := NewArray()
			for _, item := range v.([]interface{}) {
				arr.AddToEnd(item) // Добавляем элементы в массив
			}
			arrays[k] = arr

		case "list":
			lst := NewList()
			for _, item := range v.([]interface{}) {
				lst.AddToTail(item) // Добавляем элементы в список
			}
			lists[k] = lst

		case "sl":
			lst := NewSingleList()
			for _, item := range v.([]interface{}) {
				lst.AddToTail(item)
			}
			singlelists[k] = lst

		case "cbt":
			cbt := NewCBTree()
			for _, item := range v.([]interface{}) {
				cbt.Add(item)
			}
			cbtrees[k] = cbt

		case "stack":
			s := NewStack()
			for _, item := range v.([]interface{}) {
				s.Push(item)
			}
			stacks[k] = s

		case "queue":
			q := NewQueue()
			for _, item := range v.([]interface{}) {
				q.Push(item)
			}
			queues[k] = q
		case "ht":
			h := NewHashTable()
			for k1, v1 := range v.(map[string]interface{}) {
				h.Set(k1, v1) // Добавляем элементы в хеш-таблицу
			}
			hashtables[k] = h
		}
	}
}

func saveToFile(file string) {
	data := make(map[string]interface{})

	data["names"] = names

	for k, v := range arrays {
		data[k] = v.Read() // Сохраняем массивы
	}

	for k, v := range lists {
		data[k] = v.Read() // Сохраняем списки
	}

	for k, v := range singlelists {
		data[k] = v.Read()
	}

	for k, v := range hashtables {
		data[k] = v.Read() // Сохраняем хеш-таблицы
	}

	for k, v := range cbtrees {
		data[k] = v.Read() // Сохраняем полные бинарные деревья
	}

	for k, v := range queues {
		data[k] = v.Read()
	}

	for k, v := range stacks {
		data[k] = v.Read()
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err) // Обработка ошибки сериализации JSON
	}

	if err := os.WriteFile(file, jsonData, 0777); err != nil {
		log.Fatalf("Failed to write file: %v", err) // Обработка ошибки записи файла. 0777 - rwx
	}
}

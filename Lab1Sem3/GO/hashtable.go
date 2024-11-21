package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const TABLE_SIZE = 100

type Element struct {
	key   string
	value string
}

type HashMap struct {
	hashmap [TABLE_SIZE]*Element
}

// вставка нового элемента в хеш-таблицу
func (hmap *HashMap) insert(key string, value string) error {
	Element := &Element{key: key, value: value} //создается новый элемент типа Element с переданными ключем и значением
	index := HashFunc(key)
	if hmap.hashmap[index] == nil {
		hmap.hashmap[index] = Element
		fmt.Println("Элемент добавлен!")
		return nil
	} else {
		if hmap.hashmap[index].key == key {
			hmap.hashmap[index] = Element //элемент заменяется новым
			fmt.Println("Элемент обновлен!")
			return nil
		}
		index++
		for i := 0; i < TABLE_SIZE; i++ { //проходимся по таблице до конца
			if index == TABLE_SIZE {
				index = 0 //если дошли до конца, начинаем проверку с начала
			}
			if hmap.hashmap[index] == nil {
				hmap.hashmap[index] = Element //поиск ближайшего свободного индекса для вставки элемента
				fmt.Println("Элемент добавлен!")
				return nil
			}
			index++
		}
	}
	fmt.Println("Недостаточно пространства!") //если все индексы хеш-таблицы заняты, выводится сообщение об ошибке и возвращается соответствующая ошибка
	return errors.New("not enough space!")
}

// удаление элемента из хеш-таблицы
func (hmap *HashMap) delete(key string) error {
	index := HashFunc(key)

	if hmap.hashmap[index] == nil {
		fmt.Println("Объект не найден!")
		return errors.New("no item found!")
	}

	if hmap.hashmap[index].key == key {
		hmap.hashmap[index] = nil //элемент удаляется путем присвоения значения nil
		fmt.Println("Элемент удален!")
		return nil
	} else {
		index++
		for i := 0; i < TABLE_SIZE; i++ { //проходимся по таблице до конца
			if index == TABLE_SIZE {
				index = 0 //если дошли до конца, начинаем проверку с начала
			}
			if hmap.hashmap[index].key == key { //поиск элемента с данным ключом на следующих индексах, если элемент найден - он удаляется
				hmap.hashmap[index] = nil
				fmt.Println("Элемент удален!")
				return nil
			}
			index++
		}
	}
	return errors.New("not enough space in hashtable!")
}

// получение элемента из хеш-таблицы по ключу
func (hmap *HashMap) get(key string) error {
	index := HashFunc(key)
	if hmap.hashmap[index] == nil {
		fmt.Println("Объекта с таким ключом нет в хеш-таблице!")
		return errors.New("no item found!")
	}
	if hmap.hashmap[index].key == key { //проверка совпадения ключей
		fmt.Printf("Объект найден! ключ: [%s] значение: [%s]\n", hmap.hashmap[index].key, hmap.hashmap[index].value)
		return nil
	} else {
		index++
		for i := 0; i < TABLE_SIZE; i++ { //проходимся по таблице до конца
			if index == TABLE_SIZE {
				index = 0 //если дошли до конца, начинаем проверку с начала
			}
			if hmap.hashmap[index].key == key {
				fmt.Printf("Объект найден! ключ: [%s] значение: [%s]\n", hmap.hashmap[index].key, hmap.hashmap[index].value)
				return nil
			}
			index++
		}
	}
	fmt.Printf("Объект не найден!")
	return errors.New("no item found!")
}

// вывод хеш-таблицы
func (hmap *HashMap) hprint() {
	for i, item := range hmap.hashmap {
		if item == nil {
			continue
		}
		fmt.Printf("Индекс: %d. Ключ: %s | Значение: %s\n", i, item.key, item.value)
	}
}

func HashFunc(key string) int {
	sum := 0
	for _, char := range key { //проходимся по каждой букве ключа
		sum += int(char) //преобразуем букву в число
	}
	return sum % TABLE_SIZE //возвращаем хеш значение
}

func main() {
	hmap := HashMap{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nВозможные операции для хэш-таблицы:")
		fmt.Println("1. Вставить элемент в хеш-таблицу")
		fmt.Println("2. Удалить элемент из хеш-таблицы")
		fmt.Println("3. Получить элемент из хеш-таблицы")
		fmt.Println("4. Вывести хеш-таблицу")
		fmt.Println("5. Выход")
		fmt.Print("Выберите операцию: ")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			fmt.Print("Введите ключ: ")
			scanner.Scan()
			key := scanner.Text()
			fmt.Print("Введите значение: ")
			scanner.Scan()
			value := scanner.Text()
			hmap.insert(key, value)
		case "2":
			fmt.Print("Введите ключ: ")
			scanner.Scan()
			key := scanner.Text()
			hmap.delete(key)
		case "3":
			fmt.Print("Введите ключ: ")
			scanner.Scan()
			key := scanner.Text()
			hmap.get(key)
		case "4":
			hmap.hprint()
		case "5":
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор операции")
		}
	}
}

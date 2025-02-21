package main

import (
	"fmt"
	"sort"
)

// Структура для хранения информации о потомках
type Node struct {
	Children []*Node // Список указателей на потомков
	Count    int     // Количество потомков
}

func main() {
	var N int
	fmt.Scan(&N) // Считываем количество людей

	tree := NewHashTable()        // Создаем хеш-таблицу для хранения дерева
	parentTable := NewHashTable() // Создаем хеш-таблицу для хранения родителей

	// Считываем данные и строим дерево
	for i := 0; i < N; i++ {
		var child, parent string
		fmt.Scan(&child, &parent) // Считываем имя потомка и его родителя

		// Добавляем родителя и потомка в хеш-таблицу tree
		if _, exists := tree.Get(parent); !exists {
			tree.Set(parent, &Node{}) // Если родитель еще не существует, создаем новую запись
		}
		if _, exists := tree.Get(child); !exists {
			tree.Set(child, &Node{}) // Если потомок еще не существует, создаем новую запись
		}

		// Получаем указатели на родителя и потомка
		parentNode, _ := tree.Get(parent)
		childNode, _ := tree.Get(child)

		// Добавляем потомка к родителю
		parentNode.(*Node).Children = append(parentNode.(*Node).Children, childNode.(*Node))

		// Сохраняем информацию о родителе для каждого потомка в parentTable
		parentTable.Set(child, parent)
	}

	// Находим родоначальника (у кого нет родителя)
	var root string
	for person := range tree.Read() {
		if _, exists := parentTable.Get(person); !exists {
			root = person // Найденный родоначальник
			break
		}
	}

	// Вычисляем количество потомков для каждого элемента
	calculateDescendants(tree, root)

	// Получаем все элементы из хеш-таблицы tree
	elements := tree.Read()

	// Преобразуем элементы в срез для сортировки
	var keys []string
	for key := range elements {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Сортируем имена по алфавиту

	// Выводим результат
	fmt.Println("Результат:")
	for _, key := range keys {
		node := elements[key].(*Node)
		fmt.Printf("%s %d\n", key, node.Count) // Выводим имя и количество потомков
	}
}

// Функция для подсчета количества потомков
func calculateDescendants(tree *HashTable, root string) int {
	node, _ := tree.Get(root) // Получаем узел по имени
	if node.(*Node).Count != 0 {
		return node.(*Node).Count // Если количество уже посчитано, возвращаем его
	}

	count := 0
	for _, child := range node.(*Node).Children {
		childName := getKeyByValue(tree, child)            // Получаем имя потомка
		count += calculateDescendants(tree, childName) + 1 // Рекурсивно считаем потомков
	}
	node.(*Node).Count = count // Записываем количество потомков
	return count
}

// Вспомогательная функция для получения ключа по значению
func getKeyByValue(tree *HashTable, value *Node) string {
	for key, val := range tree.Read() {
		if val == value {
			return key // Возвращаем имя, соответствующее значению
		}
	}
	return "" // Если значение не найдено, возвращаем пустую строку
}

package Lists

import "fmt"

// узел односвязного или двусвязного списка.
type node struct {
	value interface{} // Значение, хранимое в узле
	next  *node       // Указатель на следующий узел
	prev  *node       // Указатель на предыдущий узел (для двусвязного списка)
}

// структура данных двусвязного списка
type List struct {
	head *node // Указатель на первый узел списка
	tail *node // Указатель на последний узел списка
}

// создает и возвращает указатель на новый пустой список
func NewList() *List {
	return &List{head: nil, tail: nil} // Инициализируем список с пустыми указателями на голову и хвост
}

// добавление нового узла в начало списка
func (l *List) AddToHead(value interface{}) {
	newNode := &node{value: value, next: l.head} // Создаем новый узел с указателем на текущий первый узел
	if l.head != nil {
		l.head.prev = newNode // Устанавливаем предыдущий узел для текущего первого узла
	} else {
		l.tail = newNode // Если список пуст, новый узел становится и первым, и последним
	}
	l.head = newNode // Новый узел становится первым узлом списка
}

// добавление нового узла в конец списка
func (l *List) AddToEnd(value interface{}) {
	newNode := &node{value: value, prev: l.tail} // Создаем новый узел с указателем на текущий последний узел
	if l.tail != nil {
		l.tail.next = newNode // Устанавливаем следующий узел для текущего последнего узла
	} else {
		l.head = newNode // Если список пуст, новый узел становится и первым, и последним
	}
	l.tail = newNode // Новый узел становится последним узлом списка
}

// удаление первого узла из списка. возвращает его значение
func (l *List) RemoveFromHead() (interface{}, error) {
	if l.head == nil {
		return nil, fmt.Errorf("List is empty") // Проверка на пустой список
	}
	value := l.head.value // Сохраняем значение первого узла
	l.head = l.head.next  // Перемещаем указатель на следующий узел
	if l.head != nil {
		l.head.prev = nil // Устанавливаем предыдущий узел для нового первого узла
	} else {
		l.tail = nil // Если список стал пустым, обнуляем указатель на последний узел
	}
	return value, nil // Возвращаем значение удаленного узла
}

// удаление последнего узла из списка. возвращает его значение
func (l *List) RemoveFromTail() (interface{}, error) {
	if l.tail == nil {
		return nil, fmt.Errorf("List is empty") // Проверка на пустой список
	}
	value := l.tail.value // Сохраняем значение последнего узла
	l.tail = l.tail.prev  // Перемещаем указатель на предыдущий узел
	if l.tail != nil {
		l.tail.next = nil // Устанавливаем следующий узел для нового последнего узла
	} else {
		l.head = nil // Если список стал пустым, обнуляем указатель на первый узел
	}
	return value, nil // Возвращаем значение удаленного узла
}

// удаление узла из списка по значению
func (l *List) RemoveByValue(value interface{}) error {
	node := l.head // Начинаем с первого узла
	for node != nil {
		if node.value == value {
			if node.prev != nil {
				node.prev.next = node.next // Обновляем указатели соседних узлов
			} else {
				l.head = node.next // Если удаляемый узел первый, обновляем указатель на первый узел
			}
			if node.next != nil {
				node.next.prev = node.prev // Обновляем указатели соседних узлов
			} else {
				l.tail = node.prev // Если удаляемый узел последний, обновляем указатель на последний узел
			}
			return nil // Узел успешно удален
		}
		node = node.next // Переходим к следующему узлу
	}
	return fmt.Errorf("value not found") // Если значение не найдено, возвращаем ошибку
}

// проверка, содержит ли список значение
func (l *List) FindByValue(value interface{}) bool {
	node := l.head // Начинаем с первого узла
	for node != nil {
		if node.value == value {
			return true // Если значение найдено, возвращаем true
		}
		node = node.next // Переходим к следующему узлу
	}
	return false // Если значение не найдено, возвращаем false
}

// возвращает копию всех значений списка
func (l *List) Read() []interface{} {
	var result []interface{} // Инициализируем срез для хранения значений
	node := l.head           // Начинаем с первого узла
	for node != nil {
		result = append(result, node.value) // Добавляем значение узла в срез
		node = node.next                    // Переходим к следующему узлу
	}
	return result // Возвращаем срез со всеми значениями
}

// структура данных односвязного списка
type SingleLinkedList struct {
	head *node // Указатель на первый узел списка
}

// создает и возвращает указатель на новый пустой односвязный список
func NewSingleList() *SingleLinkedList {
	return &SingleLinkedList{head: nil} // Инициализируем список с пустым указателем на голову
}

// добавление нового узла в начало списка
func (l *SingleLinkedList) AddToHead(value interface{}) {
	newNode := &node{value: value, next: nil} // Создаем новый узел
	if l.head == nil {
		l.head = newNode // Если список пуст, новый узел становится первым
	} else {
		newNode.next = l.head // Устанавливаем следующий узел для нового узла
		l.head = newNode      // Новый узел становится первым узлом списка
	}
}

// добавление нового узла в конец списка
func (l *SingleLinkedList) AddToEnd(value interface{}) {
	newNode := &node{value: value, next: nil} // Создаем новый узел
	if l.head == nil {
		l.head = newNode // Если список пуст, новый узел становится первым
	} else {
		current := l.head // Начинаем с первого узла
		for current.next != nil {
			current = current.next // Переходим к следующему узлу
		}
		current.next = newNode // Новый узел становится последним узлом списка
	}
}

// удаление первого узла из списка. возвращает его значение
func (l *SingleLinkedList) RemoveFromHead() (interface{}, error) {
	if l.head == nil {
		return nil, fmt.Errorf("List is empty") // Проверка на пустой список
	}
	value := l.head.value // Сохраняем значение первого узла
	l.head = l.head.next  // Перемещаем указатель на следующий узел
	return value, nil     // Возвращаем значение удаленного узла
}

// удаление последнего узла из списка. возвращает его значение
func (l *SingleLinkedList) RemoveFromTail() (interface{}, error) {
	if l.head == nil {
		return nil, fmt.Errorf("List is empty") // Проверка на пустой список
	}
	if l.head.next == nil {
		value := l.head.value // Сохраняем значение первого узла
		l.head = nil          // Если в списке только один узел, список становится пустым
		return value, nil
	}
	current := l.head // Начинаем с первого узла
	for current.next.next != nil {
		current = current.next // Переходим к предпоследнему узлу
	}
	value := current.next.value // Сохраняем значение последнего узла
	current.next = nil          // Удаляем последний узел
	return value, nil           // Возвращаем значение удаленного узла
}

// удаление узла из списка по значению
func (l *SingleLinkedList) RemoveByValue(value interface{}) error {
	if l.head.value == value {
		l.RemoveFromHead() // Если удаляемый узел первый, используем метод удаления первого узла
		return nil
	}
	previous := l.head     // Предыдущий узел
	current := l.head.next // Текущий узел
	for current != nil {
		if current.value == value {
			previous.next = current.next // Обновляем указатель предыдущего узла
			return nil                   // Узел успешно удален
		}
		previous = current // Переходим к следующему узлу
		current = current.next
	}
	return fmt.Errorf("value not found") // Если значение не найдено, возвращаем ошибку
}

// проверка, содержит ли список значение
func (l *SingleLinkedList) FindByValue(value interface{}) bool {
	node := l.head // Начинаем с первого узла
	for node != nil {
		if node.value == value {
			return true // Если значение найдено, возвращаем true
		}
		node = node.next // Переходим к следующему узлу
	}
	return false // Если значение не найдено, возвращаем false
}

// возвращает копию всех значений списка
func (l *SingleLinkedList) Read() []interface{} {
	var result []interface{} // Инициализируем срез для хранения значений
	node := l.head           // Начинаем с первого узла
	for node != nil {
		result = append(result, node.value) // Добавляем значение узла в срез
		node = node.next                    // Переходим к следующему узлу
	}
	return result // Возвращаем срез со всеми значениями
}

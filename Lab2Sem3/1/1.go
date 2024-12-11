package main

import (
	"errors"
	"fmt"
	"strings"
)

// Узел односвязного или двусвязного списка
type Node struct {
	value interface{} // Значение, хранимое в узле
	next  *Node       // Указатель на следующий узел
	prev  *Node       // Указатель на предыдущий узел (для двусвязного списка)
}

// Структура данных двусвязного списка
type List struct {
	head *Node // Указатель на первый узел списка
	tail *Node // Указатель на последний узел списка
}

// Создает и возвращает указатель на новый пустой список
func NewList() *List {
	return &List{head: nil, tail: nil}
}

// Добавление нового узла в начало списка
func (l *List) AddToHead(value interface{}) {
	newNode := &Node{value: value, next: l.head} // Создаем новый узел
	if l.head != nil {
		l.head.prev = newNode // Устанавливаем предыдущий узел для текущего первого узла
	} else {
		l.tail = newNode // Если список пуст, новый узел становится и первым, и последним
	}
	l.head = newNode // Новый узел становится первым узлом списка
}

// Удаление первого узла из списка. Возвращает его значение
func (l *List) RemoveFromHead() (interface{}, error) {
	if l.head == nil {
		return nil, errors.New("list is empty") // Проверка на пустой список
	}
	value := l.head.value // Сохраняем значение первого узла
	l.head = l.head.next  // Перемещаем указатель на следующий узел
	if l.head != nil {
		l.head.prev = nil // Устанавливаем предыдущий узел для нового первого узла
	} else {
		l.tail = nil // Если список стал пустым, обнуляем указатель на последний узел
	}
	return value, nil
}

// Возвращает копию всех значений списка
func (l *List) Read() []interface{} {
	var result []interface{} // Инициализируем срез для хранения значений
	node := l.head           // Начинаем с первого узла
	for node != nil {
		result = append(result, node.value) // Добавляем значение узла в срез
		node = node.next                    // Переходим к следующему узлу
	}
	return result // Возвращаем срез со всеми значениями
}

// Структура данных стека, реализованная с помощью двусвязного списка
type Stack struct {
	list *List // Указатель на список, используемый для хранения элементов стека
}

// Создает и возвращает указатель на новый пустой стек
func NewStack() *Stack {
	return &Stack{list: NewList()} // Создаем новый стек с пустым списком
}

// Добавление элемента в вершину стека
func (s *Stack) Push(value interface{}) {
	s.list.AddToHead(value) // Добавляем элемент в начало списка (вершину стека)
}

// Удаление элемента из вершины стека. Возвращает его значение
func (s *Stack) Pop() (interface{}, error) {
	return s.list.RemoveFromHead() // Удаляем и возвращаем элемент из начала списка (вершины стека)
}

// Возвращает копию всех элементов стека
func (s *Stack) Read() []interface{} {
	return s.list.Read() // Возвращаем копию всех элементов списка
}

// Функция для проверки и исправления XML-строки
func FixXML(xml string) (string, error) {
	var stack *Stack = NewStack()
	var fixedXML strings.Builder
	var tagStart int = -1

	for i := 0; i < len(xml); i++ {
		if xml[i] == '<' {
			tagStart = i
		} else if xml[i] == '>' && tagStart != -1 {
			tag := xml[tagStart+1 : i]
			if tag[0] != '/' { // Открывающий тег
				stack.Push(tag)
				fixedXML.WriteString(xml[tagStart : i+1])
			} else { // Закрывающий тег
				closingTag := tag[1:]
				if top, err := stack.Pop(); err == nil {
					if top != closingTag {
						// Найден несоответствующий тег, добавляем недостающие теги
						fixedXML.WriteString(fmt.Sprintf("</%s>", top))
						stack.Push(top)
						stack.Push(closingTag)
					}
					fixedXML.WriteString(xml[tagStart : i+1])
				} else {
					// Нет открывающего тега для текущего закрывающего
					fixedXML.WriteString(fmt.Sprintf("<%s>", closingTag))
					stack.Push(closingTag)
					fixedXML.WriteString(xml[tagStart : i+1])
				}
			}
			tagStart = -1
		} else {
			fixedXML.WriteByte(xml[i])
		}
	}

	// Добавляем недостающие закрывающие теги
	for _, t := range stack.Read() {
		fixedXML.WriteString(fmt.Sprintf("</%s>", t))
	}

	return fixedXML.String(), nil
}

func main() {
	xml := `<a> </b> <a> <b> <a> <b> </a> </b>`
	fixedXML, err := FixXML(xml)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Fixed XML:", fixedXML)
	}
}

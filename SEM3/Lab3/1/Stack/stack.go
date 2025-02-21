package Stack

import (
	"1/Lists"
)

// структура данных стека, реализованная с помощью двусвязного списка
type Stack struct {
	list *Lists.List // Указатель на список, используемый для хранения элементов стека
}

// создает и возвращает указатель на новый пустой стек.
func NewStack() *Stack {
	return &Stack{list: Lists.NewList()} // Создаем новый стек с пустым списком
}

// добавление элемента в вершину стека.
func (s *Stack) Push(value interface{}) {
	s.list.AddToEnd(value) // Добавляем элемент в начало списка (вершину стека)
}

// удаление элемента из вершины стека. возвращает его значение
func (s *Stack) Pop() (interface{}, error) {
	return s.list.RemoveFromTail() // Удаляем и возвращаем элемент из начала списка (вершины стека)
}

// возвращает копию всех элементов стека.
func (s *Stack) Read() []interface{} {
	return s.list.Read() // Возвращаем копию всех элементов списка
}

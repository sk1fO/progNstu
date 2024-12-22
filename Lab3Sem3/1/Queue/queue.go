package Queue

import (
	"1/Lists"
)

// структура данных очереди, реализованная с помощью двусвязного списка
type Queue struct {
	list *Lists.List // Указатель на список, используемый для хранения элементов очереди
}

// создает и возвращает указатель на новую пустую очередь
func NewQueue() *Queue {
	return &Queue{list: Lists.NewList()} // Создаем новую очередь с пустым списком
}

// добавление элемента в конец очереди
func (q *Queue) Push(value interface{}) {
	q.list.AddToEnd(value) // Добавляем элемент в конец списка
}

// удаление элемента. возвращает его значение
func (q *Queue) Pop() (interface{}, error) {
	return q.list.RemoveFromHead() // Удаляем и возвращаем элемент из начала списка
}

// возвращает копию всех элементов очереди.
func (q *Queue) Read() []interface{} {
	return q.list.Read() // Возвращаем копию всех элементов списка
}

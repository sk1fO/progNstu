package main

import "fmt"

// структура данных массива
type Array struct {
	data []interface{} // Поле для хранения элементов массива
}

// создает и возвращает новый пустой массив
func NewArray() *Array {
	return &Array{data: make([]interface{}, 0)}
}

// добавляет элемент в массив по указанному индексу
func (a *Array) AddAtIndex(index int, value interface{}) error {
	if index < 0 || index > len(a.data) {
		return fmt.Errorf("index out of bounds") // Проверка на выход за границы массива
	}
	a.data = append(a.data[:index], append([]interface{}{value}, a.data[index:]...)...)
	return nil
}

// добавляет элемент в конец массива
func (a *Array) AddToEnd(value interface{}) {
	a.data = append(a.data, value)
}

// возвращает элемент массива по указанному индексу
func (a *Array) Get(index int) (interface{}, error) {
	if index < 0 || index >= len(a.data) {
		return nil, fmt.Errorf("index out of bounds") // Проверка на выход за границы массива
	}
	return a.data[index], nil
}

// удаляет элемент из массива по указанному индексу
func (a *Array) RemoveAtIndex(index int) error {
	if index < 0 || index >= len(a.data) {
		return fmt.Errorf("index out of bounds") // Проверка на выход за границы массива
	}
	a.data = append(a.data[:index], a.data[index+1:]...)
	return nil
}

// заменяет элемент массива по указанному индексу
func (a *Array) ReplaceAtIndex(index int, value interface{}) error {
	if index < 0 || index >= len(a.data) {
		return fmt.Errorf("index out of bounds") // Проверка на выход за границы массива
	}
	a.data[index] = value
	return nil
}

// возвращает длину массива
func (a *Array) Length() int {
	return len(a.data)
}

// возвращает копию массива
func (a *Array) Read() []interface{} {
	return a.data
}

package Array

import (
	"errors"
	"unsafe"
)

const maxSize = 100 // Максимальный размер массива

// Array структура для представления массива
type Array struct {
	data   unsafe.Pointer // Указатель на данные массива
	length int            // Текущая длина массива
}

// NewArray создает новый массив
func NewArray() *Array {
	data := make([]interface{}, maxSize) // Создаем фиксированный массив с максимальным размером
	return &Array{
		data:   unsafe.Pointer(&data[0]), // Преобразуем массив в указатель
		length: 0,                        // Начальная длина массива
	}
}

// AddToEnd добавляет элемент в конец массива
func (a *Array) AddToEnd(element interface{}) error {
	if a.length >= maxSize { // Проверяем, не переполнен ли массив
		return errors.New("array overflow")
	}
	*(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(a.length)*unsafe.Sizeof(interface{}(nil)))) = element // Добавляем элемент в конец массива
	a.length++                                                                                                     // Увеличиваем длину массива
	return nil
}

// AddAtIndex добавляет элемент по указанному индексу
func (a *Array) AddAtIndex(index int, element interface{}) error {
	if index < 0 || index > a.length { // Проверяем, допустим ли индекс
		return errors.New("bad index error")
	}
	if a.length >= maxSize { // Проверяем, не переполнен ли массив
		return errors.New("array overflow")
	}
	for i := a.length; i > index; i-- { // Сдвигаем элементы вправо, начиная с конца
		*(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(i)*unsafe.Sizeof(interface{}(nil)))) = *(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(i-1)*unsafe.Sizeof(interface{}(nil))))
	}
	*(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(index)*unsafe.Sizeof(interface{}(nil)))) = element // Вставляем элемент по указанному индексу
	a.length++                                                                                                  // Увеличиваем длину массива
	return nil
}

// Get получает элемент по указанному индексу
func (a *Array) Get(index int) (interface{}, error) {
	if index < 0 || index >= a.length { // Проверяем, допустим ли индекс
		return nil, errors.New("bad index error")
	}
	return *(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(index)*unsafe.Sizeof(interface{}(nil)))), nil // Возвращаем элемент по указанному индексу
}

// RemoveAtIndex удаляет элемент по указанному индексу
func (a *Array) RemoveAtIndex(index int) error {
	if index < 0 || index >= a.length { // Проверяем, допустим ли индекс
		return errors.New("bad index error")
	}
	for i := index; i < a.length-1; i++ { // Сдвигаем элементы влево, начиная с указанного индекса
		*(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(i)*unsafe.Sizeof(interface{}(nil)))) = *(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(i+1)*unsafe.Sizeof(interface{}(nil))))
	}
	a.length-- // Уменьшаем длину массива
	return nil
}

// ReplaceAtIndex заменяет элемент по указанному индексу
func (a *Array) ReplaceAtIndex(index int, element interface{}) error {
	if index < 0 || index >= a.length { // Проверяем, допустим ли индекс
		return errors.New("bad index error")
	}
	*(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(index)*unsafe.Sizeof(interface{}(nil)))) = element // Заменяем элемент по указанному индексу
	return nil
}

// Length возвращает текущую длину массива
func (a *Array) Length() int {
	return a.length // Возвращаем текущую длину массива
}

// Read возвращает все элементы массива в виде среза
func (a *Array) Read() []interface{} {
	result := make([]interface{}, a.length) // Создаем срез для хранения элементов массива
	for i := 0; i < a.length; i++ {         // Проходим по всем элементам массива
		result[i] = *(*interface{})(unsafe.Pointer(uintptr(a.data) + uintptr(i)*unsafe.Sizeof(interface{}(nil)))) // Копируем элемент в срез
	}
	return result // Возвращаем срез с элементами массива
}

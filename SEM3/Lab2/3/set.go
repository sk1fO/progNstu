package main

// Структура множества, основанного на хеш-таблице
type Set struct {
	hashTable *HashTable
}

// Создает и возвращает указатель на новое множество
func NewSet() *Set {
	return &Set{
		hashTable: NewHashTable(),
	}
}

// Добавляет элемент в множество
func (s *Set) Add(key string) {
	s.hashTable.Set(key, true) // Используем ключ и значение true
}

// Удаляет элемент из множества
func (s *Set) Remove(key string) {
	s.hashTable.Delete(key)
}

// Проверяет наличие элемента в множестве
func (s *Set) Contains(key string) bool {
	_, exists := s.hashTable.Get(key)
	return exists
}

// Возвращает все элементы множества в виде среза строк
func (s *Set) Read() []string {
	elements := s.hashTable.Read()
	result := make([]string, 0, len(elements))
	for key := range elements {
		result = append(result, key)
	}
	return result
}

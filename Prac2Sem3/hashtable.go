package main

import (
	"math"
)

const initialCapacity = 16 // Начальная емкость хеш-таблицы
const loadFactor = 0.75    // Коэффициент загрузки

// Структура данных хеш-таблицы
type HashTable struct {
	buckets  []*bucketList // Массив указателей на списки узлов (бакеты)
	size     int           // Текущий размер хеш-таблицы (количество элементов)
	capacity int           // Текущая емкость хеш-таблицы
}

// Список узлов в бакете
type bucketList struct {
	head *bucket // Голова списка
	tail *bucket // Хвост списка
}

// Узел хеш-таблицы, содержащий ключ и значение.
type bucket struct {
	key   string      // Ключ
	value interface{} // Значение
	next  *bucket     // Указатель на следующий узел
}

// Создает и возвращает указатель на новую хеш-таблицу с начальной емкостью.
func NewHashTable() *HashTable {
	return &HashTable{
		buckets:  make([]*bucketList, initialCapacity), // Инициализируем массив бакетов с начальной емкостью
		size:     0,                                    // Начальный размер хеш-таблицы равен 0
		capacity: initialCapacity,                      // Начальная емкость
	}
}

// Вычисляет хеш-значение для ключа и возвращает индекс в массиве бакетов
func (h *HashTable) hash(key string) int {
	return int(math.Abs(float64(hashString(key) % h.capacity))) // Вычисляем хеш-значение и берем остаток от деления на текущую емкость
}

// Хеш-функция для строки
func hashString(s string) int {
	var hash int = 0 // Инициализируем хеш-значение
	for i := 0; i < len(s); i++ {
		hash = (hash << 5) ^ int(s[i]) // Обновляем хеш-значение с использованием битовых операций
	}
	return hash
}

// Добавление или обновление пары ключ-значение в хеш-таблице
func (h *HashTable) Set(key string, value interface{}) {
	if float64(h.size)/float64(h.capacity) >= loadFactor {
		h.rehash()
	}

	index := h.hash(key) // Вычисляем индекс для ключа
	if h.buckets[index] == nil {
		h.buckets[index] = &bucketList{}
	}

	list := h.buckets[index]
	current := list.head
	for current != nil {
		if current.key == key {
			current.value = value // Если ключ уже существует, обновляем значение
			return
		}
		current = current.next
	}

	// Если ключа нет, добавляем новый узел в начало списка
	newBucket := &bucket{key: key, value: value}
	newBucket.next = list.head
	list.head = newBucket
	if list.tail == nil {
		list.tail = newBucket
	}
	h.size++ // Увеличиваем размер хеш-таблицы
}

// Перехеширование хеш-таблицы
func (h *HashTable) rehash() {
	oldBuckets := h.buckets
	h.capacity *= 2 // Увеличиваем емкость вдвое
	h.buckets = make([]*bucketList, h.capacity)
	h.size = 0

	for _, list := range oldBuckets {
		if list != nil {
			current := list.head
			for current != nil {
				h.Set(current.key, current.value) // Переносим элементы в новую хеш-таблицу
				current = current.next
			}
		}
	}
}

// Возвращает значение по ключу, если ключ существует, иначе возвращает nil и false.
func (h *HashTable) Get(key string) (interface{}, bool) {
	index := h.hash(key) // Вычисляем индекс для ключа
	if h.buckets[index] == nil {
		return nil, false
	}

	list := h.buckets[index]
	current := list.head
	for current != nil {
		if current.key == key {
			return current.value, true // Возвращаем значение и true, если ключ найден
		}
		current = current.next
	}
	return nil, false // Возвращаем nil и false, если ключ не найден
}

// Удаляет пару ключ-значение из хеш-таблицы по ключу
func (h *HashTable) Delete(key string) {
	index := h.hash(key) // Вычисляем индекс для ключа
	if h.buckets[index] == nil {
		return
	}

	list := h.buckets[index]
	if list.head == nil {
		return
	}

	if list.head.key == key {
		list.head = list.head.next
		if list.head == nil {
			list.tail = nil
		}
		h.size--
		return
	}

	prev := list.head
	current := list.head.next
	for current != nil {
		if current.key == key {
			prev.next = current.next
			if current.next == nil {
				list.tail = prev
			}
			h.size--
			return
		}
		prev = current
		current = current.next
	}
}

// Возвращает копию всех элементов хеш-таблицы в виде среза строк
func (h *HashTable) Read() map[string]interface{} {
	result := make(map[string]interface{}) // Инициализируем срез для хранения результатов
	for _, list := range h.buckets {
		if list != nil {
			current := list.head
			for current != nil {
				result[current.key] = current.value // Формируем строку "ключ:значение" и добавляем в срез
				current = current.next
			}
		}
	}
	return result // Возвращаем срез со всеми элементами
}

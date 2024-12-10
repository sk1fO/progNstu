package main

const initialCapacity = 16 // Начальная емкость хеш-таблицы

// структура данных хеш-таблицы
type HashTable struct {
	buckets []*bucket // Массив указателей на узлы (бакеты)
	size    int       // Текущий размер хеш-таблицы (количество элементов)
}

// узел хеш-таблицы, содержащий ключ и значение.
type bucket struct {
	key   string      // Ключ
	value interface{} // Значение
}

// создает и возвращает указатель на новую хеш-таблицу с начальной емкостью.
func NewHashTable() *HashTable {
	return &HashTable{
		buckets: make([]*bucket, initialCapacity), // Инициализируем массив бакетов с начальной емкостью
		size:    0,                                // Начальный размер хеш-таблицы равен 0
	}
}

// вычисляет хеш-значение для ключа и возвращает индекс в массиве бакетов
func (h *HashTable) hash(key string) int {
	return hashString(key) % len(h.buckets) // Вычисляем хеш-значение и берем остаток от деления на длину массива бакетов
}

// хеш-функция для строки
func hashString(s string) int {
	var hash int = 0 // Инициализируем хеш-значение
	for i := 0; i < len(s); i++ {
		hash = (hash << 5) ^ int(s[i]) // Обновляем хеш-значение с использованием битовых операций
	}
	return hash
}

// добавление или обновление пары ключ-значение в хеш-таблице
func (h *HashTable) Set(key string, value interface{}) {
	index := h.hash(key) // Вычисляем индекс для ключа
	for _, b := range h.buckets[index:] {
		if b != nil && b.key == key {
			b.value = value // Если ключ уже существует, обновляем значение
			return
		}
	}
	h.buckets[index] = &bucket{key: key, value: value} // Если ключа нет, добавляем новый узел
	h.size++                                           // Увеличиваем размер хеш-таблицы
}

// возвращает значение по ключу, если ключ существует, иначе возвращает nil и false.
func (h *HashTable) Get(key string) (interface{}, bool) {
	index := h.hash(key) // Вычисляем индекс для ключа
	for _, b := range h.buckets[index:] {
		if b != nil && b.key == key {
			return b.value, true // Возвращаем значение и true, если ключ найден
		}
	}
	return nil, false // Возвращаем nil и false, если ключ не найден
}

// удаляет пару ключ-значение из хеш-таблицы по ключу
func (h *HashTable) Delete(key string) {
	index := h.hash(key) // Вычисляем индекс для ключа
	for i, b := range h.buckets[index:] {
		if b != nil && b.key == key {
			h.buckets[index+i] = nil // Удаляем узел
			h.size--                 // Уменьшаем размер хеш-таблицы
			return
		}
	}
}

// возвращает копию всех элементов хеш-таблицы в виде среза строк
func (h *HashTable) Read() []interface{} {
	var result []interface{} // Инициализируем срез для хранения результатов
	for _, b := range h.buckets {
		if b != nil {
			result = append(result, b.key+":"+b.value.(string)) // Формируем строку "ключ:значение" и добавляем в срез
		}
	}
	return result // Возвращаем срез со всеми элементами
}

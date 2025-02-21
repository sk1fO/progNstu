package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Функция для вычисления пересечения двух множеств
func intersection(set1, set2 *Set) int {
	count := 0
	elements := set1.Read()
	for _, elem := range elements {
		if set2.Contains(elem) {
			count++
		}
	}
	return count
}

// Функция для решения задачи
func findMaxIntersection(subsets []*Set) (set1, set2 *Set, maxIntersection int) {
	maxIntersection = -1
	for i := 0; i < len(subsets); i++ {
		for j := i + 1; j < len(subsets); j++ {
			intersectionCount := intersection(subsets[i], subsets[j])
			if intersectionCount > maxIntersection {
				maxIntersection = intersectionCount
				set1, set2 = subsets[i], subsets[j]
			}
		}
	}
	return set1, set2, maxIntersection
}

func main() {
	fmt.Print("Введите множества: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	//[[1, 2, 3],[2, 3, 4, 6],[5, 6],[3, 4, 5, 6]]
	// Парсим входные данные
	var subsets [][]int
	err := json.Unmarshal([]byte(input), &subsets)
	if err != nil {
		fmt.Printf("Ошибка при парсинге входных данных: %v\n", err)
		return
	}

	// Преобразуем входные данные в множества
	var sets []*Set
	for _, subset := range subsets {
		set := NewSet()
		for _, elem := range subset {
			set.Add(fmt.Sprintf("%d", elem)) // Добавляем элементы как строки
		}
		sets = append(sets, set)
	}

	// Находим пару с максимальным пересечением
	set1, set2, maxIntersection := findMaxIntersection(sets)

	// Выводим результат
	fmt.Printf("Пара множеств с максимальным пересечением: %v и %v\n", set1.Read(), set2.Read())
	fmt.Printf("Количество общих элементов: %d\n", maxIntersection)
}

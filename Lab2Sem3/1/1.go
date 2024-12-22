package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// checkTags проверяет, является ли переданная строка корректным XML-подобным тегом.
func checkTags(xml string) bool {
	// Проверяем, что строка начинается с '<' и заканчивается '>'.
	if xml[0] != '<' || xml[len(xml)-1] != '>' {
		return false
	}

	// Создаем стек для хранения открывающих тегов.
	stack := NewStack()
	i := 0

	// Основной цикл для обработки строки.
	for {
		// Если достигли конца строки, выходим из цикла.
		if i == len(xml) {
			break
		}

		// Проверяем, что текущий символ - '<'.
		if xml[i] != '<' {
			return false
		}
		i++

		// Флаг для определения, является ли тег закрывающим.
		closingTag := false
		if i < len(xml) && xml[i] == '/' {
			closingTag = true
			i++
		}

		// Проверяем, что следующий символ - строчная буква.
		if i >= len(xml) || !unicode.IsLower(rune(xml[i])) {
			return false
		}

		// Считываем тег.
		tag := string(xml[i])
		i++
		for i < len(xml) && unicode.IsLower(rune(xml[i])) {
			tag += string(xml[i])
			i++
		}

		// Проверяем, что тег заканчивается '>'.
		if i >= len(xml) || xml[i] != '>' {
			return false
		}
		i++

		// Если тег открывающий, добавляем его в стек.
		if !closingTag {
			stack.Push(tag)
		} else {
			// Если тег закрывающий, проверяем его соответствие последнему открывающему.
			lastTag, err := stack.Pop()
			if err != nil {
				return false
			}
			if lastTag != tag {
				return false
			}
		}
	}

	// Если стек пуст, значит все теги закрыты корректно.
	return len(stack.Read()) == 0
}

// changeSymbol изменяет один символ в строке и проверяет, является ли новая строка корректной.
func changeSymbol(str []rune, j rune, result *[][]rune) {
	for i := 0; i < len(str); i++ {
		// Создаем копию строки.
		temp := make([]rune, len(str))
		copy(temp, str)

		// Изменяем символ.
		temp[i] = j

		// Проверяем, является ли новая строка корректной.
		if checkTags(string(temp)) {
			*result = append(*result, temp)
		}
	}
}

func main() {
	// Приглашение пользователю ввести строку.
	fmt.Print("Введите XML строку: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()
	str = strings.ReplaceAll(str, " ", "")

	// Преобразуем строку в массив рун для удобства работы с символами.
	stringRunes := []rune(str)

	// Счетчики для символов '<', '>', '/' и букв.
	openedChevron := 0
	closedChevron := 0
	slashes := 0
	letters := make(map[rune]int)

	// Результирующий массив для корректных строк.
	result := [][]rune{}

	// Подсчитываем количество символов '<', '>', '/' и букв.
	for _, r := range stringRunes {
		if r == '<' {
			openedChevron++
		} else if r == '>' {
			closedChevron++
		} else if r == '/' {
			slashes++
		}
		letters[r]++
	}

	// Перебираем возможные символы для замены.
	for _, j := range []rune("<>/qwertyuiopasdfghjklzxcvbnm") {
		// Проверяем условия для замены символов.
		if j == '<' && openedChevron%2 != 0 {
			changeSymbol(stringRunes, j, &result)
		} else if j == '>' && closedChevron%2 != 0 {
			changeSymbol(stringRunes, j, &result)
		} else if j == '/' && closedChevron/2 != slashes {
			changeSymbol(stringRunes, j, &result)
		} else if count, exists := letters[j]; exists && count%2 != 0 {
			changeSymbol(stringRunes, j, &result)
		}
	}

	// Выводим все найденные корректные строки.
	for _, res := range result {
		fmt.Println("Исправленная строка:", string(res))
	}
}

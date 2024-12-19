package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Функция для проверки корректности XML-строки
func isCorrectXML(xml string) bool {
	stack := []string{}
	tags := strings.Split(xml, ">")
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		tag = strings.TrimSuffix(tag, "<")
		if !strings.HasPrefix(tag, "/") {
			// Это открывающий тег
			stack = append(stack, tag)
		} else {
			// Это закрывающий тег
			if len(stack) == 0 {
				return false
			}
			lastTag := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if lastTag != tag[1:] {
				return false
			}
		}
	}
	return len(stack) == 0
}

// Функция для исправления XML-строки
func fixXML(xml string) string {
	for i := 0; i < len(xml); i++ {
		// Попробуем заменить каждый символ
		for c := 'a'; c <= 'z'; c++ {
			if string(xml[i]) == string(c) {
				continue
			}
			fixedXML := xml[:i] + string(c) + xml[i+1:]
			if isCorrectXML(fixedXML) {
				return fixedXML
			}
		}
		if xml[i] == '<' {
			fixedXML := xml[:i] + ">" + xml[i+1:]
			if isCorrectXML(fixedXML) {
				return fixedXML
			}
		}
		if xml[i] == '>' {
			fixedXML := xml[:i] + "<" + xml[i+1:]
			if isCorrectXML(fixedXML) {
				return fixedXML
			}
		}
		if xml[i] == '/' {
			fixedXML := xml[:i] + "" + xml[i+1:]
			if isCorrectXML(fixedXML) {
				return fixedXML
			}
		}
	}

	// Если простые замены не помогли, попробуем более сложные варианты
	for i := 0; i < len(xml); i++ {
		if xml[i] == '<' {
			// Попробуем заменить на закрывающий тег
			for j := i + 1; j < len(xml); j++ {
				if xml[j] == '>' {
					// Найдем имя тега
					tagName := xml[i+1 : j]
					if tagName != "" && !strings.HasPrefix(tagName, "/") {
						// Создадим закрывающий тег
						closingTag := "</" + tagName + ">"
						fixedXML := xml[:i] + closingTag + xml[j+1:]
						if isCorrectXML(fixedXML) {
							return fixedXML
						}
					}
				}
			}
		}
		if xml[i] == '/' {
			// Попробуем удалить слэш
			fixedXML := xml[:i] + xml[i+1:]
			if isCorrectXML(fixedXML) {
				return fixedXML
			}
		}
	}

	return xml // Если ничего не помогло, вернем исходную строку
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	xml := scanner.Text()

	fixedXML := fixXML(xml)
	fmt.Println(fixedXML)
}

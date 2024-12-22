package main

import (
	"fmt"
	"strings"
)

func isOpeningTag(s string) bool {
	return len(s) >= 2 && s[0] == '<' && s[1] != '/' && s[len(s)-1] == '>'
}

func isClosingTag(s string) bool {
	return len(s) >= 3 && s[0] == '<' && s[1] == '/' && s[len(s)-1] == '>'
}

func getTagName(s string) string {
	if isOpeningTag(s) {
		return s[1 : len(s)-1]
	} else if isClosingTag(s) {
		return s[2 : len(s)-1]
	}
	return ""
}

func fixTag(tag string, expectedTagName string) string {
	// Try to fix the tag by changing one character
	if isClosingTag(tag) {
		// If it's a closing tag, try to fix it to match the expected opening tag
		return "</" + expectedTagName + ">"
	} else if isOpeningTag(tag) {
		// If it's an opening tag, try to fix it to match the expected closing tag
		return "<" + expectedTagName + ">"
	} else {
		// If it's neither, try to fix it to be a closing tag
		return "</" + expectedTagName + ">"
	}
}

func restoreXML(input string) string {
	tags := strings.Fields(input)
	stack := []string{}

	for _, tag := range tags {
		if isOpeningTag(tag) {
			stack = append(stack, tag)
		} else if isClosingTag(tag) {
			if len(stack) == 0 {
				continue
			}
			openingTag := stack[len(stack)-1]
			openingTagName := getTagName(openingTag)
			closingTagName := getTagName(tag)

			if openingTagName != closingTagName {
				// Try to fix the closing tag
				fixedTag := fixTag(tag, openingTagName)
				if isClosingTag(fixedTag) && getTagName(fixedTag) == openingTagName {
					tag = fixedTag
				}
			}

			if openingTagName == closingTagName {
				stack = stack[:len(stack)-1]
			}
		} else {
			// If the tag is invalid, try to fix it
			if len(stack) > 0 {
				openingTag := stack[len(stack)-1]
				openingTagName := getTagName(openingTag)
				fixedTag := fixTag(tag, openingTagName)
				if isClosingTag(fixedTag) && getTagName(fixedTag) == openingTagName {
					tag = fixedTag
					stack = stack[:len(stack)-1]
				}
			}
		}
	}

	// Reconstruct the XML string
	var result strings.Builder
	for _, tag := range stack {
		result.WriteString(tag + " ")
	}

	return strings.TrimSpace(result.String())
}

func main() {
	input := "<a> //a> <a> <b> </b> </a>"
	fmt.Println("Исходная строка:", input)
	restored := restoreXML(input)
	fmt.Println("Восстановленная строка:", restored)
}

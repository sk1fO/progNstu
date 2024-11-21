package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var Alphabet = make(map[rune]int)
var AlphabetAlpha = make(map[int]rune)
var N int

func init() {
	rusAlphabet := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	engAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbolAlphabet := " ,./';[]!@#$%^&*()-_=+{}<>?№:1234567890`~" + `|\"`

	AlphabetString := rusAlphabet + engAlphabet + symbolAlphabet
	AlphabetRune := []rune(AlphabetString)

	for number, symbol := range AlphabetRune {
		Alphabet[symbol] = number
		AlphabetAlpha[number] = symbol
	}

	N = len(Alphabet)
}

func main() {

	clearConsole()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите текст для шифрования: ")
	scanner.Scan()
	input := scanner.Text()

	fmt.Print("Введите ключ: ")
	scanner.Scan()
	keyInput := scanner.Text()

	runeKey := []rune(keyInput)
	runeInput := []rune(input)

	crypted := encryption(runeInput, runeKey)
	fmt.Println(crypted)

	runeCrypted := []rune(crypted)
	decrypted := decryption(runeCrypted, runeKey)
	fmt.Println(decrypted)
}

// функция шифрования
func encryption(rawText, key []rune) string {

	// удлиннение ключа до размера входной строки
	for len(key) < len(rawText) {
		key = append(key, key...)
	}

	result := "" // объявление строки - результата
	var symbol rune
	for i, sym := range rawText {
		char, ok := Alphabet[sym] // проверяем наличие символа в алфавите
		if ok {
			symbol = AlphabetAlpha[(char+Alphabet[key[i]])%N] // если есть, то Zi = Xi + Ki (mod N)
			result += string(symbol)
		} else {
			result += string(sym) // иначе оставляем "как есть"
		}

	}

	return result // возврат результата
}

// функция расшифровки
func decryption(rawText, key []rune) string {

	//удлиняем ключ до размера входной строки
	for len(key) < len(rawText) {
		key = append(key, key...)
	}

	result := "" // переменная результата
	var symbol rune
	for i, sym := range rawText {

		char, ok := Alphabet[sym] // если символ есть в алфавите
		if ok {
			symbol = AlphabetAlpha[(char-Alphabet[key[i]]+N)%N] // Xi = Zi - Ki + N (mod N)
			result += string(symbol)
		} else {
			result += string(sym) // если нет, оставляем "как есть"
		}

	}

	return result

}

// Linux очистка консоли
func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

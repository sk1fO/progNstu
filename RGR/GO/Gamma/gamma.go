package Gamma

var Alphabet = make(map[rune]int)
var AlphabetAlpha = make(map[int]rune)
var N int

func init() {
	rusAlphabet := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	engAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbolAlphabet := "	" + `|\"`

	AlphabetString := rusAlphabet + engAlphabet + symbolAlphabet
	AlphabetRune := []rune(AlphabetString)

	for number, symbol := range AlphabetRune {
		Alphabet[symbol] = number
		AlphabetAlpha[number] = symbol
	}

	N = len(Alphabet)
}

// функция шифрования, принимает исходный текст и ключ в виде строки
// возвращает шифротекст текст в виде строки
func Encrypt(input, keyInput string) (cipherText string) {

	key := []rune(keyInput)
	rawText := []rune(input)

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

// функция расшифровки, принимает шифротекст и ключ в виде строки
// возвращает расшифрованый текст в виде строки
func Decrypt(ciferText, keyInput string) (plainText string) {

	rawText := []rune(ciferText)
	key := []rune(keyInput)
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

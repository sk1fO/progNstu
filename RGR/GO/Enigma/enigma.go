package Enigma

import (
	"strings"
)

var N int

var reflect = make(map[rune]rune)

// заполнение глобальной переменной reflect и N
func init_reflector() {
	Alphabet := "TDLGЕЖKRНИUЛЪБФQДJВVEРАЦЗЮXХFЯNОYУOBПЩZHЬЭСIAWШЫSMТPЙКГЧCМ" // алфавит отражателя
	AlphabetAlpha := "МCЧГКЙPТMSЫШWAIСЭЬHZЩПBOУYОNЯFХXЮЗЦАРEVВJДQФБЪЛUИНRKЖЕGLDT"
	AlphabetRune := []rune(Alphabet)
	AlphabetAlphaRune := []rune(AlphabetAlpha)

	N = len(AlphabetRune)

	for i := range N {
		reflect[AlphabetAlphaRune[i]] = AlphabetRune[i]
	}
}

// Функция описывающая первый ротор, принимает: символ на замену, инверсия?, сдвиг
func first_rotor(symbol rune, inv bool, rotation int) rune {
	var rotor = make(map[rune]rune)

	Alphabet := "АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯABCDEFGHIJKLMNOPQRSTUVWXYZ"      // алфавит входной
	AlphabetAlpha := "HKЯFLРЙCАNZVБDНТIКMAЭЗГДЩЬSХBЛЮUQЦYМШЪЖRСGФОЧOXPУTЫEИПWЕJВ" // алфавит сдвига
	AlphabetRune := []rune(Alphabet)                                              // перевод в срез рун
	AlphabetAlphaRune := []rune(AlphabetAlpha)

	// заполнение map прямой или обратной, в зависимости до или после рефлектора
	if inv {
		for i := range N {
			rotor[AlphabetAlphaRune[(i+rotation)%N]] = AlphabetRune[i]
		}
	} else {
		for i := range N {
			rotor[AlphabetRune[i]] = AlphabetAlphaRune[(i+rotation)%N]
		}
	}

	// Если есть символ в алфавите, то возвращаем его замену
	// иначе возвращаем "как есть"
	result, ok := rotor[symbol]
	if ok {
		return result
	} else {
		return symbol
	}
}

// Функция описывающая второй ротор, принимает: символ на замену, инверсия?, сдвиг
func second_rotor(symbol rune, inv bool, rotation int) rune {
	var rotor = make(map[rune]rune)

	Alphabet := "ТZCЛЫAЭOЖМKYVИXБФШUHПJЦОЕАGВDЪЩSЯTЬЗСРNХIPWЙЧLBRДКГFУEЮНMQ"      // алфавит входной
	AlphabetAlpha := "СLDIBHEPВNFПOYЬЯОТНУKХMЮФWTАSЙЩКЕГМШAДЪЧQЖRCVXЗGРZБJUИЦЫЛЭ" // алфавит сдвига
	AlphabetRune := []rune(Alphabet)                                              // перевод в срез рун
	AlphabetAlphaRune := []rune(AlphabetAlpha)

	if inv {
		for i := range N {
			rotor[AlphabetAlphaRune[(i+rotation)%N]] = AlphabetRune[i]
		}
	} else {
		for i := range N {
			rotor[AlphabetRune[i]] = AlphabetAlphaRune[(i+rotation)%N]
		}
	}

	// Если есть символ в алфавите, то возвращаем его замену
	// иначе возвращаем "как есть"
	result, ok := rotor[symbol]
	if ok {
		return result
	} else {
		return symbol
	}
}

// Функция описывающая третий ротор, принимает: символ на замену, инверсия?, сдвиг
func third_rotor(symbol rune, inv bool, rotation int) rune {
	var rotor = make(map[rune]rune)

	Alphabet := "СLDIBHEPВNFПOYЬЯОТНУKХMЮФWTАSЙЩКЕГМШAДЪЧQЖRCVXЗGРZБJUИЦЫЛЭ"      // алфавит входной
	AlphabetAlpha := "ZЯКРУQЛGЭХYUHSWMПНГACОЮIФBODTRТЬEБДXЙLPЩFЪВЕVШЧKNЗЖJМЫАИСЦ" // алфавит сдвига
	AlphabetRune := []rune(Alphabet)                                              // перевод в срез рун
	AlphabetAlphaRune := []rune(AlphabetAlpha)

	if inv {
		for i := range N {
			rotor[AlphabetAlphaRune[(i+rotation)%N]] = AlphabetRune[i]
		}
	} else {
		for i := range N {
			rotor[AlphabetRune[i]] = AlphabetAlphaRune[(i+rotation)%N]
		}
	}

	// Если есть символ в алфавите, то возвращаем его замену
	// иначе возвращаем "как есть"
	result, ok := rotor[symbol]
	if ok {
		return result
	} else {
		return symbol
	}
}

// Функция описывающая работу отражателя принимает: символ для отражения
func reflector(symbol rune) rune {

	result, ok := reflect[symbol]
	if ok {
		return result // если символ есть в алфавите, то возвращаем его замену
	} else {
		return symbol // иначе возвращаем символ "как есть"
	}
}

// Функция шифрования
// Принимает: входной текст строка и массив стартовых положений роторов
// Возвращает: зашифрованный текст
func Encrypt(input string, startpos [3]int) string {
	init_reflector()               // инициализация отражателя
	var char rune                  // инициализация символа для перевода
	result := ""                   // итоговая строка
	input = strings.ToUpper(input) // приведение к верхнему регистру
	runeInput := []rune(input)     // перевод в срез рун

	pos1 := startpos[0] // чтение стартовых позиций
	pos2 := startpos[1]
	pos3 := startpos[2]

	for _, sym := range runeInput {
		// каждый символ проходит через 3 ротора
		char = first_rotor(sym, false, pos1)
		char = second_rotor(char, false, pos2)
		char = third_rotor(char, false, pos3)
		char = reflector(char)
		char = third_rotor(char, true, pos3)
		char = second_rotor(char, true, pos2)
		char = first_rotor(char, true, pos1)

		result += string(char)

		// механизм сдвига роторов
		if pos1 < N-1 {
			pos1++
		} else {
			pos2++
			pos1 = 0
		}

		if pos2 == N {
			pos3++
			pos2 = 0
		}

		if pos3 == N {
			pos3 = 0
		}
	}
	return result
}

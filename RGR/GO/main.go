package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"GO/AES128"
	"GO/Enigma"
	"GO/Gamma"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	f, _ := os.Open("passhash.bin")
	passScaner := bufio.NewScanner(f)
	passScaner.Scan()
	passhash := passScaner.Bytes()
	f.Close()

	for {
		clearConsole()
		fmt.Print("Введиите пароль: ")

		scanner.Scan()
		s := scanner.Text()

		h := sha1.New()
		h.Write([]byte(s))
		bs := h.Sum(nil)

		//os.WriteFile("passhash.bin", bs, os.FileMode(0777))

		if string(bs) == string(passhash) {
			break

		} else {
			fmt.Print("Неверный пароль! (enter, чтобы продолжить)")
			scanner.Scan()
		}
	}

out:
	for {
		clearConsole()
		fmt.Println("Меню\n<1> Шифрование AES128\n<2> Шифрование Гамма\n<3> Шифрование машиной Энигма\n<4> Смена пароля\n<5> Выход")
		fmt.Print(">>> ")
		scanner.Scan()
		menuPos := scanner.Text()

		switch menuPos {
		case "1":
			clearConsole()

			fmt.Println("Шифрование AES128\n<1> Шифрование\n<2> Дешифрование")
			fmt.Print(">>> ")
			scanner.Scan()

			aesMenu := scanner.Text()
			switch aesMenu {
			case "1":
				clearConsole()
				fmt.Println("Шифрование AES128\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":

				case "2":
					clearConsole()
					fmt.Print("Шифрование AES128\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите ключ: ")
					scanner.Scan()
					inputKey := scanner.Text()

					cipher := AES128.Encrypt(input, inputKey)
					fmt.Println("Зашифрованное сообщение:", cipher)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						permissions := 0777 // or whatever you need
						byteArray := []byte(cipher)
						os.WriteFile("file.txt", byteArray, os.FileMode(permissions))
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}

				}

			case "2":
				clearConsole()
				fmt.Println("Шифрование AES128\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					clearConsole()
					f, _ := os.Open("file.txt")
					Fscanner := bufio.NewScanner(f)
					Fscanner.Scan()
					input := string(Fscanner.Bytes())
					fmt.Println("Шифрование AES128\nТекст для дешифрования: ", input)
					fmt.Print("Введите ключ: ")
					scanner.Scan()
					inputKey := scanner.Text()

					plainText := AES128.Decrypt(input, inputKey)
					fmt.Println("Расшифрованное сообщение:", plainText)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование AES128\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите ключ: ")
					scanner.Scan()
					inputKey := scanner.Text()

					plainText := AES128.Decrypt(input, inputKey)
					fmt.Println("Расшифрованное сообщение:", plainText)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}
				}
			}

		case "2":
			clearConsole()

			fmt.Println("Шифрование Гамма\n<1> Шифрование\n<2> Дешифрование")
			fmt.Print(">>> ")
			scanner.Scan()

			gammaMenu := scanner.Text()
			switch gammaMenu {
			case "1":
				clearConsole()
				fmt.Println("Шифрование Гамма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					// file
				case "2":
					clearConsole()
					fmt.Print("Шифрование Гамма\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите гамму: ")
					scanner.Scan()
					inputKey := scanner.Text()

					cipher := Gamma.Encrypt(input, inputKey)
					fmt.Println("Зашифрованное сообщение:", cipher)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}

				}

			case "2":
				clearConsole()
				fmt.Println("Шифрование Гамма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					// file
				case "2":
					clearConsole()
					fmt.Print("Шифрование Гамма\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите гамму: ")
					scanner.Scan()
					inputKey := scanner.Text()

					plainText := Gamma.Decrypt(input, inputKey)
					fmt.Println("Расшифрованное сообщение:", plainText)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}
				}
			}

		case "3":
			clearConsole()

			fmt.Println("Шифрование машиной Энигма\n<1> Шифрование\n<2> Дешифрование")
			fmt.Print(">>> ")
			scanner.Scan()

			aesMenu := scanner.Text()
			switch aesMenu {
			case "1":
				clearConsole()
				fmt.Println("Шифрование машиной Энигма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					// file
				case "2":
					clearConsole()
					fmt.Print("Шифрование машиной Энигма\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите начальные положения роторов (3): ")
					scanner.Scan()
					inputKey := scanner.Text()

					var rotorPos [3]int
					inputKeySlice := strings.Split(inputKey, " ")
					for i := range 3 {
						rotorPos[i], _ = strconv.Atoi(inputKeySlice[i])
					}

					cipher := Enigma.Encrypt(input, rotorPos)
					fmt.Println("Зашифрованное сообщение:", cipher)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}
				}

			case "2":
				clearConsole()
				fmt.Println("Шифрование машиной Энигма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					// file
				case "2":
					clearConsole()
					fmt.Print("Шифрование машиной Энигма\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					fmt.Print("Введите начальные положения роторов (3): ")
					scanner.Scan()
					inputKey := scanner.Text()

					var rotorPos [3]int
					inputKeySlice := strings.Split(inputKey, " ")
					for i := range 3 {
						rotorPos[i], _ = strconv.Atoi(inputKeySlice[i])
					}

					plainText := Enigma.Encrypt(input, rotorPos)
					fmt.Println("Расшифрованное сообщение:", plainText)
					fmt.Print("\nЗаписать в файл? (y/n): ")
					scanner.Scan()
					fileYN := scanner.Text()
					switch fileYN {
					case "y":
						fmt.Println("Файл записан")
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					case "n":
						fmt.Println("Нажмите enter для возврата в меню")
						scanner.Scan()
					}
				}
			}
		case "4":
			clearConsole()

			f, _ := os.Open("passhash.bin")
			passScaner := bufio.NewScanner(f)
			passScaner.Scan()
			passhash := passScaner.Bytes()
			f.Close()

			fmt.Print("Введиите текущий пароль: ")

			scanner.Scan()
			s := scanner.Text()

			h := sha1.New()
			h.Write([]byte(s))
			bs := h.Sum(nil)

			if string(bs) == string(passhash) {
				fmt.Print("Введиите новый пароль: ")

				scanner.Scan()
				s := scanner.Text()

				h := sha1.New()
				h.Write([]byte(s))
				bs := h.Sum(nil)
				os.WriteFile("passhash.bin", bs, os.FileMode(0777))

			} else {
				fmt.Println("Неверный пароль!")
				scanner.Scan()
			}
		case "5":
			break out
		}
	}
}

// Linux очистка консоли
func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

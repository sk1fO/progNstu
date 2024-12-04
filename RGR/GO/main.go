package main

import (
	"bufio"       // работа с чтением/записью
	"crypto/sha1" // хеш-функция для входа по паролю
	"fmt"         // стандартная библиотека
	"os"          // работа с ос и файлами
	"os/exec"     // работа с терминалом

	"GO/AES128" // подключение шифров
	"GO/Enigma"
	"GO/Gamma"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) // объявляем сканер, для чтения файлов
	f, _ := os.Open("passhash.bin")       // открываем бинарный файл с хеш-суммой пароля
	passScaner := bufio.NewScanner(f)     // еще один сканер, но уже для файла
	passScaner.Scan()                     // читаем данные файла
	passhash := passScaner.Bytes()        //присваиваем их переменной в виде среза байт
	f.Close()

	for { // бесконечный цикл ожидания ввода верного пароля, без него не войти в меню
		clearConsole()
		fmt.Print("Введиите пароль: ")

		scanner.Scan()
		s := scanner.Text()

		h := sha1.New()    // новый объект для хеш-функции sha1
		h.Write([]byte(s)) // вычисляем хеш
		bs := h.Sum(nil)   //суммируем

		// приводим к строке, чтобы можно было сравнить
		// срезы байт сравнивать нельзя
		if string(bs) == string(passhash) {
			break // всё хорошо - останавливаем цикл, переходим в меню
		} else {
			fmt.Print("Неверный пароль! (enter, чтобы продолжить)")
			scanner.Scan()
		}
	}

out: // нужен, чтобы иметь возможность выйти из беск. цикла
	for {
		// печатаем меню и вызываем подменю для каждого шифра
		clearConsole()
		fmt.Println("Меню\n<1> Шифрование AES128\n<2> Шифрование Гамма\n<3> Шифрование машиной Энигма\n<4> Смена пароля\n<5> Выход")
		fmt.Print(">>> ")
		scanner.Scan()
		menuPos := scanner.Text()

		switch menuPos {
		case "1":
			AES128Menu()
		case "2":
			GammaMenu()
		case "3":
			EnigmaMachineMenu()
		case "4":
			clearConsole()
			// повторяем процедуру из начала
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

				h := sha1.New() //вычисляем хеш сумму нового пароля
				h.Write([]byte(s))
				bs := h.Sum(nil)
				os.WriteFile("passhash.bin", bs, os.FileMode(0777)) // перезаписываем бинарный файл

			} else {
				fmt.Println("Неверный пароль!")
				scanner.Scan()
			}
		case "5":
			break out //останавливаем цикл - выход из программы
			// можно заменить на return
		}
	}
}

// Linux очистка консоли
func clearConsole() {
	cmd := exec.Command("clear") // выполняем команду
	cmd.Stdout = os.Stdout       // в консоли
	cmd.Run()
}

// Функция чтения файла
func fileRead(name string) (textFromFile string) {
	path := "input/" + name
	f, ok := os.Open(path)
	if ok != nil { // обработка ошибок
		fmt.Println("ошибка чтения файла")
		return "" // возвращаем пустую строку
	}
	defer f.Close() // во время выхода из ф-ции файл будет закрыт

	fs := bufio.NewScanner(f)
	fs.Scan()

	return string(fs.Bytes())
}

// Функция записи в файл
func writeFile(name, text string) {
	byteArray := []byte(text)
	path := "output/" + name
	ok := os.WriteFile(path, byteArray, os.FileMode(0777))
	if ok != nil { // обработка ошибок
		fmt.Println("ошибка записи в файл")
		return
	}

	fmt.Println("Файл записан")
}

// функция подменю для шифра энигма
func EnigmaMachineMenu() {
outMenuEnigma:
	for {
		scanner := bufio.NewScanner(os.Stdin)
		clearConsole()

		fmt.Println("Шифрование машиной Энигма\n<1> Шифрование\n<2> Дешифрование")
		fmt.Print(">>> ")
		scanner.Scan()

		enigmaMenu := scanner.Text()
		switch enigmaMenu {
		case "1":
			for {
				clearConsole()
				fmt.Println("Шифрование машиной Энигма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				enigmaMenuFile := scanner.Text()
				switch enigmaMenuFile {
				case "1":
					clearConsole()

					inputFileName := ""
					for inputFileName == "" {
						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					var rotorPos [3]int
					a, b, c := 0, 0, 0

					fmt.Print("Введите начальные положения роторов(3): ")
					_, err := fmt.Scan(&a, &b, &c)
					for err != nil {
						fmt.Println("Некорректные положения роторов(3)")
						_, err = fmt.Scan(&a, &b, &c)
					}

					rotorPos[0] = a
					rotorPos[1] = b
					rotorPos[2] = c

					cipher := Enigma.Encrypt(input, rotorPos)

					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}
							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование машиной Энигма\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					var rotorPos [3]int
					a, b, c := 0, 0, 0
					fmt.Print("Введите начальные положения роторов(3): ")
					_, err := fmt.Scan(&a, &b, &c)
					for err != nil {
						fmt.Println("Некорректные положения роторов (3)")
						_, err = fmt.Scan(&a, &b, &c)
					}

					rotorPos[0] = a
					rotorPos[1] = b
					rotorPos[2] = c

					cipher := Enigma.Encrypt(input, rotorPos)

					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}
							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}

					}

				default:
					fmt.Println("Повторите попытку")
					fmt.Println("Нажмите enter, чтобы продолжить")
					scanner.Scan()
				}
			}

		case "2":
			for {
				clearConsole()
				fmt.Println("Шифрование машиной Энигма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					clearConsole()

					inputFileName := ""

					for inputFileName == "" {
						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					var rotorPos [3]int
					a, b, c := 0, 0, 0
					fmt.Print("Введите начальные положения роторов(3): ")
					_, err := fmt.Scan(&a, &b, &c)
					for err != nil {
						fmt.Println("Некорректные положения роторов(3)")
						_, err = fmt.Scan(&a, &b, &c)
					}

					rotorPos[0] = a
					rotorPos[1] = b
					rotorPos[2] = c

					plainText := Enigma.Encrypt(input, rotorPos)

					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование машиной Энигма\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					var rotorPos [3]int
					a, b, c := 0, 0, 0
					fmt.Print("Введите начальные положения роторов(3): ")
					_, err := fmt.Scan(&a, &b, &c)
					for err != nil {
						fmt.Println("Некорректные положения роторов(3)")
						_, err = fmt.Scan(&a, &b, &c)
					}

					rotorPos[0] = a
					rotorPos[1] = b
					rotorPos[2] = c

					plainText := Enigma.Encrypt(input, rotorPos)

					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuEnigma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}
				}
			}

		default:
			fmt.Println("Повторите попытку")
			fmt.Println("Нажмите enter, чтобы продолжить")
			scanner.Scan()
		}
	}
}

// функция подменю для шифра гаммой
func GammaMenu() {
	scanner := bufio.NewScanner(os.Stdin)

outMenuGamma:
	for {
		clearConsole()

		fmt.Println("Шифрование Гамма\n<1> Шифрование\n<2> Дешифрование")
		fmt.Print(">>> ")
		scanner.Scan()

		gammaMenu := scanner.Text()
		switch gammaMenu {
		case "1":
			for {
				clearConsole()
				fmt.Println("Шифрование Гамма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				gammaMenuFile := scanner.Text()
				switch gammaMenuFile {
				case "1":
					clearConsole()
					inputFileName := ""

					for inputFileName == "" {

						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					inputGamma := ""
					for inputGamma == "" {
						fmt.Print("Введите гамму: ")
						scanner.Scan()
						inputGamma = scanner.Text()
						if inputGamma == "" {
							fmt.Println("Некорректная гамма")
						}
					}

					cipher := Gamma.Encrypt(input, inputGamma)
					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						case "n":
							clearConsole()
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}

					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование Гамма\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					inputGamma := ""
					for inputGamma == "" {
						fmt.Print("Введите гамму: ")
						scanner.Scan()
						inputGamma = scanner.Text()
						if inputGamma == "" {
							fmt.Println("Некорректная гамма")
						}
					}

					cipher := Gamma.Encrypt(input, inputGamma)
					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}
							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}

					}
				default:
					fmt.Println("Повторите попытку")
					fmt.Println("Нажмите enter, чтобы продолжить")
					scanner.Scan()
				}
			}

		case "2":
			for {
				clearConsole()
				fmt.Println("Шифрование Гамма\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				gammaMenuFile := scanner.Text()
				switch gammaMenuFile {
				case "1":
					clearConsole()

					inputFileName := ""

					for inputFileName == "" {
						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					inputGamma := ""
					for inputGamma == "" {
						fmt.Print("Введите гамму: ")
						scanner.Scan()
						inputGamma = scanner.Text()
						if inputGamma == "" {
							fmt.Println("Некорректная гамма")
						}
					}

					plainText := Gamma.Decrypt(input, inputGamma)
					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""
							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						default:
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma
						}
					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование Гамма\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					inputGamma := ""
					for inputGamma == "" {
						fmt.Print("Введите гамму: ")
						scanner.Scan()
						inputGamma = scanner.Text()
						if inputGamma == "" {
							fmt.Println("Некорректная гамма")
						}
					}

					plainText := Gamma.Decrypt(input, inputGamma)

					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuGamma

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}
				}
			}
		default:
			fmt.Println("Повторите попытку")
			fmt.Println("Нажмите enter, чтобы продолжить")
			scanner.Scan()
		}
	}
}

// функция подменю для шифра AES
func AES128Menu() {
outMenuAes:
	for {
		scanner := bufio.NewScanner(os.Stdin)
		clearConsole()

		fmt.Println("Шифрование AES128\n<1> Шифрование\n<2> Дешифрование")
		fmt.Print(">>> ")
		scanner.Scan()

		aesMenu := scanner.Text()
		switch aesMenu {
		case "1":
			for {
				clearConsole()
				fmt.Println("Шифрование AES128\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					clearConsole()
					inputFileName := ""

					for inputFileName == "" {

						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					inputKey := ""
					for inputKey == "" {
						fmt.Print("Введите ключ: ")
						scanner.Scan()
						inputKey = scanner.Text()
						if inputKey == "" {
							fmt.Println("Некорректный ключ")
						}
					}

					cipher := AES128.Encrypt(input, inputKey)
					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()

						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}
							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						case "n":
							clearConsole()
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование AES128\nВведите текст для шифрования: ")
					scanner.Scan()
					input := scanner.Text()

					inputKey := ""
					for inputKey == "" {
						fmt.Print("Введите ключ: ")
						scanner.Scan()
						inputKey = scanner.Text()
						if inputKey == "" {
							fmt.Println("Некорректный ключ")
						}
					}

					cipher := AES128.Encrypt(input, inputKey)
					for {
						fmt.Println("Зашифрованное сообщение:", cipher)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()
							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}
							writeFile(outputFileName, cipher)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}

				default:
					fmt.Println("Повторите попытку")
					fmt.Println("Нажмите enter, чтобы продолжить")
					scanner.Scan()
				}
			}

		case "2":
			for {
				clearConsole()
				fmt.Println("Шифрование AES128\n<1> Ввод исходного текста из файла\n<2> Ввод исходного текста через консоль")
				fmt.Print(">>> ")
				scanner.Scan()

				aesMenuFile := scanner.Text()
				switch aesMenuFile {
				case "1":
					clearConsole()

					inputFileName := ""

					for inputFileName == "" {
						fmt.Print("Введите название файла из директории input: ")
						scanner.Scan()
						inputFileName = scanner.Text()
						if inputFileName == "" {
							fmt.Println("Название не может быть пустым")
						}
					}

					input := fileRead(inputFileName)
					fmt.Println(input)

					inputKey := ""
					for inputKey == "" {
						fmt.Print("Введите ключ: ")
						scanner.Scan()
						inputKey = scanner.Text()
						if inputKey == "" {
							fmt.Println("Некорректный ключ")
						}
					}

					plainText := AES128.Decrypt(input, inputKey)
					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes
						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}

					}

				case "2":
					clearConsole()
					fmt.Print("Шифрование AES128\nВведите текст для дешифрования: ")
					scanner.Scan()
					input := scanner.Text()

					inputKey := ""
					for inputKey == "" {
						fmt.Print("Введите ключ: ")
						scanner.Scan()
						inputKey = scanner.Text()
						if inputKey == "" {
							fmt.Println("Некорректный ключ")
						}
					}

					plainText := AES128.Decrypt(input, inputKey)

					for {
						fmt.Println("Расшифрованное сообщение:", plainText)
						fmt.Print("\nЗаписать в файл? (y/n): ")
						scanner.Scan()
						fileYN := scanner.Text()
						switch fileYN {
						case "y":
							clearConsole()

							outputFileName := ""

							for outputFileName == "" {
								fmt.Print("Введите название файла для записи: ")
								scanner.Scan()
								outputFileName = scanner.Text()
								if outputFileName == "" {
									fmt.Println("Название не может быть пустым")
								}
							}

							writeFile(outputFileName, plainText)

							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						case "n":
							fmt.Println("Нажмите enter для возврата в меню")
							scanner.Scan()
							break outMenuAes

						default:
							fmt.Println("Повторите попытку")
							fmt.Println("Нажмите enter, чтобы продолжить")
							scanner.Scan()
						}
					}
				}
			}

		default:
			fmt.Println("Повторите попытку")
			fmt.Println("Нажмите enter, чтобы продолжить")
			scanner.Scan()
		}
	}
}

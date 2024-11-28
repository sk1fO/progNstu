package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"GO/AES128"
)

func main() {
	clearConsole()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите текст для шифрования: ")
	scanner.Scan()
	input := scanner.Text()

	fmt.Print("Введите ключ: ")
	scanner.Scan()
	inputKey := scanner.Text()

	cipher := AES128.Encrypt(input, inputKey)
	fmt.Println(cipher)

	plain := AES128.Decrypt(cipher, inputKey)
	fmt.Println(plain)

}

// Linux очистка консоли
func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите текст для шифрования: ")
	scanner.Scan()
	input := scanner.Text()

	byteInput := []byte(input)

	var plainText [16]byte

	for len(byteInput) != 0 {
		if len(byteInput) < 16 {
			byteInput = append(byteInput, 0)
		} else {
			plainText = [16]byte(byteInput)
			fmt.Println(plainText)
			byteInput = byteInput[16:]
		}

	}

	//fmt.Println(byteInput)
}

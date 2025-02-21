package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Block структура для представления блока
type Block struct {
	Width  int
	Height int
}

// Сортировка блоков по ширине с использованием сортировки вставками
func sortBlocks(array *Array) {
	for i := 1; i < array.Length(); i++ {
		current, _ := array.Get(i)
		currentBlock := current.(Block)

		j := i - 1
		for j >= 0 {
			prev, _ := array.Get(j)
			prevBlock := prev.(Block)

			if prevBlock.Width > currentBlock.Width {
				// Сдвигаем блок вправо
				array.ReplaceAtIndex(j+1, prevBlock)
				j--
			} else {
				break
			}
		}
		// Вставляем текущий блок на правильное место
		array.ReplaceAtIndex(j+1, currentBlock)
	}
}

func main() {
	fmt.Print("Введите кол-во блоков: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	blockNum, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalln("неверное кол-во блоков", input)
	}

	blocks := NewArray()
	for i := 0; i < blockNum; i++ {
		var w, h int
		fmt.Scan(&w, &h)
		blocks.AddToEnd(Block{w, h})
	}

	// Сортируем блоки по ширине
	sortBlocks(blocks)

	// Инициализируем массив dp
	dp := make([]int, blocks.Length())
	for i := 0; i < blocks.Length(); i++ {
		block, _ := blocks.Get(i)
		dp[i] = block.(Block).Height
	}

	// Построение пирамиды
	for i := 1; i < blocks.Length(); i++ {
		for j := 0; j < i; j++ {
			blockI, _ := blocks.Get(i)
			blockJ, _ := blocks.Get(j)
			if blockJ.(Block).Width < blockI.(Block).Width {
				if dp[j]+blockI.(Block).Height > dp[i] {
					dp[i] = dp[j] + blockI.(Block).Height
				}
			}
		}
	}

	// Находим максимальную высоту
	maxHeight := 0
	for _, height := range dp {
		if height > maxHeight {
			maxHeight = height
		}
	}

	fmt.Println("Максимальная высота:", maxHeight)
}

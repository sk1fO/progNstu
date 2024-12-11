package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Block struct {
	width  int
	height int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	blocks := make([]Block, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		parts := strings.Split(scanner.Text(), " ")
		w, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])
		blocks[i] = Block{width: w, height: h}
	}

	// Сортируем блоки по убыванию ширины. Если ширина одинаковая, сортируем по высоте.
	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].width == blocks[j].width {
			return blocks[i].height > blocks[j].height
		}
		return blocks[i].width > blocks[j].width
	})

	// Используем динамическое программирование для поиска LIS по высоте.
	dp := []int{}
	for _, block := range blocks {
		pos := sort.Search(len(dp), func(i int) bool { return dp[i] >= block.height })
		if pos == len(dp) {
			dp = append(dp, block.height)
		} else {
			dp[pos] = block.height
		}
	}

	// Максимальная высота пирамиды равна сумме высот в LIS.
	maxHeight := 0
	for _, h := range dp {
		maxHeight += h
	}

	fmt.Println(maxHeight)
}

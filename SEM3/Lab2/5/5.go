package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// TreeNode представляет узел бинарного дерева поиска
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Insert добавляет элемент в бинарное дерево поиска
func (root *TreeNode) Insert(val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}
	if val < root.Val {
		root.Left = root.Left.Insert(val)
	} else {
		root.Right = root.Right.Insert(val)
	}
	return root
}

// FindNodesWithTwoChildren находит все узлы с двумя детьми
func (root *TreeNode) FindNodesWithTwoChildren() []int {
	var result []int
	if root == nil {
		return result
	}

	// Рекурсивно проверяем левое поддерево
	result = append(result, root.Left.FindNodesWithTwoChildren()...)

	// Если у текущего узла есть оба ребенка, добавляем его значение в список
	if root.Left != nil && root.Right != nil {
		result = append(result, root.Val)
	}

	// Рекурсивно проверяем правое поддерево
	result = append(result, root.Right.FindNodesWithTwoChildren()...)

	return result
}

func main() {
	// Ввод данных
	fmt.Print("Введите числа: ") // 8 3 10 1 6 14 4 7 13

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	numSlice := strings.Split(scanner.Text(), " ")

	nums := make([]int, 0, 10)
	for _, i := range numSlice {
		num, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println("Некорректное число", i)
		}
		nums = append(nums, num)
	}

	// Построение дерева
	var BST *TreeNode
	for _, num := range nums {
		if num == 0 {
			break
		}
		BST = BST.Insert(num)
	}

	// Нахождение узлов с двумя детьми
	nodes := BST.FindNodesWithTwoChildren()

	// Сортируем результат, так как обход в глубину может не дать отсортированный список
	sort.Ints(nodes)

	// Вывод результата
	fmt.Println(nodes)
}

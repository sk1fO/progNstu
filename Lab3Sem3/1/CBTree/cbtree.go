package CBTree

// структура узла дерева
type TreeNode struct {
	Value interface{}
	Left  *TreeNode
	Right *TreeNode
}

// структура данных полного бинарного дерева
type CBTree struct {
	root *TreeNode // корень дерева
	size int       // размер дерева
}

// создает и возвращает указатель на новое пустое полное бинарное дерево
func NewCBTree() *CBTree {
	return &CBTree{root: nil, size: 0}
}

// добавляет элемент в полное бинарное дерево
func (t *CBTree) Add(value interface{}) {
	t.size++
	if t.root == nil {
		t.root = &TreeNode{Value: value}
		return
	}
	queue := []*TreeNode{t.root}
	for {
		node := queue[0]
		queue = queue[1:]
		if node.Left == nil {
			node.Left = &TreeNode{Value: value}
			return
		}
		if node.Right == nil {
			node.Right = &TreeNode{Value: value}
			return
		}
		queue = append(queue, node.Left, node.Right)
	}
}

// проверяет, содержит ли дерево заданное значение.
func (t *CBTree) Find(value interface{}) bool {
	queue := []*TreeNode{t.root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.Value == value {
			return true
		}
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return false
}

// проверяет, является ли дерево полным бинарным деревом.
func (t *CBTree) IsComplete() bool {
	if t.root == nil {
		return true
	}
	queue := []*TreeNode{t.root}
	isLeafExpected := false
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if isLeafExpected && (node.Left != nil || node.Right != nil) {
			return false
		}

		if node.Left == nil && node.Right != nil {
			return false
		}

		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		} else {
			isLeafExpected = true
		}
	}
	return true
}

// возвращает копию всех элементов дерева.
func (t *CBTree) Read() []interface{} {
	if t.root == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0, t.size)
	queue := []*TreeNode{t.root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node.Value)
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return result
}

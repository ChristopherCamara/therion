package main

type TreeNode[T comparable] struct {
	data     T
	children []*TreeNode[T]
}

func (t *TreeNode[T]) addChild(data T) *TreeNode[T] {
	child := new(TreeNode[T])
	child.data = data
	t.children = append(t.children, child)
	return child
}

func (t *TreeNode[T]) findNode(data T) (*TreeNode[T], bool) {
	var found *TreeNode[T]
	queue := []*TreeNode[T]{t}
	for len(queue) > 0 {
		queue = append(queue, t.children...)
		if t.data == data {
			return found, true
		}
	}
	return nil, false
}

func NewTree[T comparable](data T) *TreeNode[T] {
	root := new(TreeNode[T])
	root.data = data
	return root
}

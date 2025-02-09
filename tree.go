package main

import (
	"time"

	"github.com/google/uuid"
)

type parent[T any] struct {
	id   uuid.UUID
	data T
}

type move[T any] struct {
	timeStamp time.Time
	parent    uuid.UUID
	data      T
	child     uuid.UUID
}

type moveLog[T any] struct {
	timeStamp time.Time
	oldParent *parent[T]
	newParent uuid.UUID
	data      T
	child     uuid.UUID
}

type treeNode[T any] struct {
	id       uuid.UUID
	data     T
	children []*treeNode[T]
}

func newTreeNode[T any]() *treeNode[T] {
	return new(treeNode[T])
}

type Tree[T any] struct {
	root *treeNode[T]
	log  []moveLog[T]
}

func getParent[T any](t *Tree[T], findChild uuid.UUID) *parent[T] {
	queue := []*treeNode[T]{t.root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range node.children {
			if child.id == findChild {
				return &parent[T]{id: node.id, data: node.data}
			}
			queue = append(queue, child)
		}
	}
	return nil
}

func ancestor() {}

func (t *Tree[T]) doOp(m move[T]) {
	t.log = append(t.log, moveLog[T]{timeStamp: m.timeStamp, oldParent: getParent[T](t, m.child), newParent: m.parent, data: m.data, child: m.child})
	queue := []*treeNode[T]{t.root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range node.children {
			queue = append(queue, child)
		}
		if node.id == m.parent {
			newChild := new(treeNode[T])
			newChild.id = m.child
			newChild.data = m.data
			node.children = append(node.children, newChild)
			queue = nil
		}
	}
}

func undoOp() {}

func redoOp() {}

func applyOp() {}

func newTree[T any]() *Tree[T] {
	tree := new(Tree[T])
	root := newTreeNode[T]()
	tree.root = root
	return tree
}

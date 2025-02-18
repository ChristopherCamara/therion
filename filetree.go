package main

import (
	"time"

	"github.com/google/uuid"
)

type file struct {
	path string
	data []byte
}

type move struct {
	timeStamp time.Time
	parent    *uuid.UUID
	file      *file
	child     *uuid.UUID
}

type log struct {
	timeStamp time.Time
	oldParent *TreeNode
	newParent *uuid.UUID
	file      *file
	child     *uuid.UUID
}

type TreeNode struct {
	parent uuid.UUID
	file   *file
}

type Tree struct {
	nodes    map[uuid.UUID]TreeNode
	children map[uuid.UUID][]uuid.UUID
}

func (t *Tree) getParent(child *uuid.UUID) *TreeNode {
	if node, ok := t.nodes[*child]; ok {
		if parent, ok := t.nodes[node.parent]; ok {
			return &parent
		}
	}
	return nil
}

func (t *Tree) isAncenstor(parent, child *uuid.UUID) bool {
	node, ok := t.nodes[*child]
	for ok {
		if node.parent == *parent {
			return true
		}
		node, ok = t.nodes[node.parent]
	}
	return false
}

func (t *Tree) doOp(m move) log {
	log := log{timeStamp: m.timeStamp, oldParent: t.getParent(m.child), newParent: m.parent, file: m.file, child: m.child}
	if m.parent != m.child && !t.isAncenstor(m.parent, m.child) {
		t.nodes[*m.child] = TreeNode{parent: *m.parent, file: m.file}
		children, ok := t.children[*m.parent]
		if ok {
			children = append(children, *m.child)
		} else {
			children = []uuid.UUID{*m.child}
		}
		t.children[*m.parent] = children
	}
	return log
}

func undoOp() {}

func redoOp() {}

func applyOp() {}

func (t *Tree) getIDFromPath(path string) *uuid.UUID {
	for id, node := range t.nodes {
		if node.file.path == path {
			return &id
		}
	}
	return nil
}

type FileTree struct {
	tree Tree
	logs []log
}

func (ft FileTree) doMove(m move) {
	log := ft.tree.doOp(m)
	ft.logs = append(ft.logs, log)
}

func (ft FileTree) AddDirectory(path string) {
}

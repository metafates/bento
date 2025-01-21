package treewidget

import "github.com/metafates/bento/textwidget"

type Item[T comparable] struct {
	id       T
	text     textwidget.Text
	children []Item[T]
}

func NewItem[T comparable](id T, text textwidget.Text, children ...Item[T]) Item[T] {
	return Item[T]{
		id:       id,
		text:     text,
		children: children,
	}
}

func (i *Item[T]) ID() T {
	return i.id
}

func (i *Item[T]) Children() []Item[T] {
	return i.children
}

func (i *Item[T]) Child(index int) (Item[T], bool) {
	if index >= len(i.children) {
		return Item[T]{}, false
	}

	return i.children[index], true
}

func (i *Item[T]) Height() int {
	return i.text.Height()
}

func (i *Item[T]) AddChild(child Item[T]) {
	i.children = append(i.children, child)
}

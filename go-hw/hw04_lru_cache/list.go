package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) insertAfter(node, newNode *ListItem) *ListItem {
	newNode.Prev = node
	if node.Next == nil {
		newNode.Next = nil
		l.back = newNode
	} else {
		newNode.Next = node.Next
		node.Next.Prev = newNode
	}
	node.Next = newNode
	l.len++
	return newNode
}

func (l *list) insertBefore(node, newNode *ListItem) *ListItem {
	newNode.Next = node
	if node.Prev == nil {
		newNode.Prev = nil
		l.front = newNode
	} else {
		newNode.Prev = node.Prev
		node.Prev.Next = newNode
	}
	node.Prev = newNode
	l.len++
	return newNode
}

func (l *list) insertBeginning(newNode *ListItem) *ListItem {
	if l.front == nil {
		l.front = newNode
		l.back = newNode
		newNode.Prev = nil
		newNode.Next = nil
		l.len++
		return newNode
	}
	return l.insertBefore(l.front, newNode)
}

func (l *list) insertEnd(newNode *ListItem) *ListItem {
	if l.back == nil {
		return l.insertBeginning(newNode)
	}
	return l.insertAfter(l.back, newNode)
}

func (l *list) PushBack(v interface{}) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.back,
	}

	return l.insertEnd(newNode)
}

func (l *list) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  l.front,
		Next:  nil,
	}

	return l.insertBeginning(newNode)
}

func (l *list) Remove(node *ListItem) {
	if node.Prev == nil {
		l.front = node.Next
	} else {
		node.Prev.Next = node.Next
	}
	if node.Next == nil {
		l.back = node.Prev
	} else {
		node.Next.Prev = node.Prev
	}
	l.len--
}

func (l *list) MoveToFront(node *ListItem) {
	l.Remove(node)
	l.PushFront(node.Value)
}

func NewList() List {
	return &list{}
}

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
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length    int
	frontItem *ListItem
	backItem  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() {
		l.length++
	}()

	newNode := &ListItem{Value: v}

	if l.frontItem == nil {
		l.frontItem = newNode
		l.backItem = newNode
	} else {
		newNode.Next = l.frontItem
		l.frontItem.Prev = newNode
		l.frontItem = newNode
	}

	return l.frontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer func() {
		l.length++
	}()

	newNode := new(ListItem)
	newNode.Value = v

	if l.backItem == nil {
		l.frontItem = newNode
		l.backItem = newNode
	} else {
		newNode.Prev = l.backItem
		l.backItem.Next = newNode
		l.backItem = newNode
	}

	return l.backItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil || l.length == 0 {
		return
	}

	if i.Next == nil && i.Prev == nil {
		l.frontItem = nil
		l.backItem = nil
		l.length--
		return
	}

	if i.Prev == nil {
		l.frontItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.backItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}

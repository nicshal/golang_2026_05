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

/*
root - метеданные списка:
- root.Next  - указывает на первый элемент списка
- root.Prev  - указывает на последний элемент списка
- root.Value - содержит длину списка.
*/
type list struct {
	root ListItem
}

func (l *list) init() List {
	l.root.Next = nil
	l.root.Prev = nil
	l.root.Value = 0
	return l
}

func NewList() List {
	return new(list).init()
}

func (l *list) Len() int { return l.root.Value.(int) }

func (l *list) Front() *ListItem {
	return l.root.Next
}

func (l *list) Back() *ListItem {
	return l.root.Prev
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := ListItem{v, nil, nil}
	if l.root.Value.(int) == 0 {
		l.root.Prev = &i
	} else {
		i.Next = l.root.Next
		l.root.Next.Prev = &i
	}
	l.root.Next = &i
	l.root.Value = l.root.Value.(int) + 1

	return l.root.Next
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := ListItem{v, nil, nil}
	if l.root.Value.(int) == 0 {
		l.root.Next = &i
	} else {
		i.Prev = l.root.Prev
		l.root.Prev.Next = &i
	}
	l.root.Prev = &i
	l.root.Value = l.root.Value.(int) + 1

	return l.root.Prev
}

func (l *list) Remove(i *ListItem) {
	// если пришел nil - ничего не делаем, возвращаемся
	if i == nil {
		return
	}

	// если удаляем последний элемент списка - просто инициализизуем список заново
	if l.root.Value.(int) == 1 {
		l.init()
		return
	}

	switch {
	case i.Next == nil: // если удаляем элемент из конца списка
		l.root.Prev = i.Prev
		i.Prev.Next = nil
	case i.Prev == nil: // если удаляем элемент из начала списка
		l.root.Next = i.Next
		i.Next.Prev = nil
	default: // если удаляем элемент из середины списка
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Next = nil
	i.Prev = nil
	l.root.Value = l.root.Value.(int) - 1
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil { // передвигаем, если элемент уже не стоит в начале списка
		i.Prev.Next = i.Next
		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.root.Prev = i.Prev
		}
		i.Prev = nil
		l.root.Next.Prev = i
		i.Next = l.root.Next
		l.root.Next = i
	}
}

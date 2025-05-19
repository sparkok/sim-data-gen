package entity

import (
	"container/list"
)

type LastNState struct {
	maxNumber int
	elements  *list.List
}

func NewLastNState(n int) *LastNState {
	return &LastNState{
		maxNumber: n,
		elements:  list.New(),
	}
}
func (t *LastNState) Add(element string) {
	if t.elements.Len() >= t.maxNumber {
		t.elements.Remove(t.elements.Back())
	}
	t.elements.PushFront(element)
}
func (t *LastNState) GetAll() []string {
	result := make([]string, t.elements.Len())
	i := 0
	for e := t.elements.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			result[i] = str
		}
		i++
	}
	return result
}

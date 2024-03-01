package main

import "fmt"

type CustomList[T any] struct {
	list  []T
	count int
}

func NewList[T any](size int) CustomList[T] {
	return CustomList[T]{
		list:  make([]T, size),
		count: 0,
	}
}
func (c *CustomList[T]) add(element T) {
	c.list = append(c.list, element)
	c.count++
}
func (c *CustomList[T]) remove(pos int) bool {
	if pos >= 0 && pos < c.count {
		c.list = append(c.list[:pos], c.list[pos+1:]...)
		c.count--
		return true
	}
	return false
}
func getZero[T any]() T {
	var result T
	return result
}
func (c *CustomList[T]) getIndex(index int) (T, error) {
	if index >= 0 && index < c.count {
		return c.list[index], fmt.Errorf("NO ELEMENT FOUND")
	}
	return getZero[T](), fmt.Errorf("NO ELEMENT FOUND")
}

package main

import (
	"container/list"
)

func main() {
	mylist := list.New()
	mylist.PushBack(1)
	mylist.PushFront(2)

}

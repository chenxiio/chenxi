package main

import (
	"fmt"
)

type UserJob struct {
	MATERIAL_LIST []string
}

func main() {
	uj := UserJob{
		MATERIAL_LIST: []string{"a", "b", "c", "d", "e"},
	}
	i := 2

	beforeI := uj.MATERIAL_LIST[:i]
	afterI := uj.MATERIAL_LIST[i:]

	fmt.Println("Before i:", beforeI)
	fmt.Println("After i:", afterI)
}

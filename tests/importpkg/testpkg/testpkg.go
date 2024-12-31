package testpkg

import "fmt"

func init() {
	fmt.Println("testpkg init")
}

func add(a, b int) int {
	return a + b
}

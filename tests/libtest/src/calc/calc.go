package main

import (
	"fmt"

	"github.com/chenxiio/chenxi/tests/libtest/src/calc/fibonacci"
)

func main() {
	var res int64
	var err error

	res, err = fibonacci.Fibonacci(30)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Result:", res)
	}

	res, err = fibonacci.Fibonacci_r(30)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Result:", res)
	}
}

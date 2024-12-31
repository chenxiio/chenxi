package fibonacci

import (
	"errors"

	"github.com/chenxiio/chenxi/tests/libtest/src/simplemath"
)

func Fibonacci(n int64) (int64, error) {
	if n < 1 {
		err := errors.New("Should be greater than 0!")
		return 0, err
	} else if n > 92 {
		err := errors.New("Should be less than 93!")
		return 0, err
	}

	var res int64 = 0
	var tmp int64 = 1
	var idx int64 = 0

	for ; idx < n; idx++ {
		res = int64(simplemath.Add(res, tmp))
		res, tmp = tmp, res
	}

	return res, nil
}

func Fibonacci_r(n int64) (int64, error) {
	if n < 1 {
		err := errors.New("Should be greater than 0!")
		return 0, err
	} else if n < 3 {
		return 1, nil
	} else if n > 92 {
		err := errors.New("Should be less than 93!")
		return 0, err
	}

	lhs, _ := Fibonacci_r(n - 1)
	rhs, _ := Fibonacci_r(n - 2)
	ret := simplemath.Add(lhs, rhs)
	return ret, nil
}

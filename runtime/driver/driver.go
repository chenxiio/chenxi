package driver

// #include <stdio.h>
// int call_function(void* func, int a, int b) {
//     int (*add)(int, int);
//     add = (int (*)(int, int))func;
//     return add(a, b);
// }
import "C"

type Driver struct {
}

func (d *Driver) Add(a, b int) int {
	return a + b
}

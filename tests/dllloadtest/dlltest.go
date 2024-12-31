package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/chenxiio/chenxi/comm/dllload"
)

// #include <stdio.h>
// int call_function(void* func, int a, int b) {
//     int (*add)(int, int);
//     add = (int (*)(int, int))func;
//     return add(a, b);
// }
import "C"

func main() {
	// 加载动态库
	lib := dllload.LoadLibrary("./test")
	//libdll := syscall.MustLoadDLL("test.dll")

	// if err != nil {
	// 	fmt.Println("LoadLibrary error:", err)
	// 	return

	defer dllload.FreeLibrary(lib)
	// 获取函数地址
	proc := dllload.GetProcAddress(lib, "add")

	var result int
	start := time.Now()
	for i := 0; i < 10000000; i++ {
		result = int((C.call_function(unsafe.Pointer(proc), C.int(1), C.int(2))))
	}

	fmt.Printf("1 + 2 = %d\n", result)
	fmt.Printf("Call time: %v\n", time.Since(start))
	start = time.Now()
	// 转换函数类型
	// fn := func() int {
	// 	ret, _, _ := syscall.SyscallN(uintptr(proc), 1, 2)
	// 	return int(ret)
	// }

	// 调用函数
	for i := 0; i < 100000000; i++ {
		result = add(1, 2)
	}

	fmt.Println("add:", result)
	fmt.Printf("Call time: %v\n", time.Since(start))
}
func add(a, b int) int {
	return a + b
}

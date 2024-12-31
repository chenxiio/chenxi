//go:build linux
// +build linux

package main

/*
 #cgo LDFLAGS: -L . -ldl -lstdc++
 #cgo CFLAGS: -I ./
#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>

int call_function(void* func, int a, int b) {
    int (*add)(int, int);
    add = (int (*)(int, int))func;
    return add(a, b);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 打开动态库
	handle := C.dlopen(C.CString("./test.so"), C.RTLD_LAZY)
	if handle == nil {
		panic("无法打开动态库")
	}
	//defer C.dlclose(handle)
	// 获取add函数
	addSymbol := C.CString("add")
	defer C.free(unsafe.Pointer(handle))
	addFuncPtr := C.dlsym(handle, addSymbol)
	if addFuncPtr == nil {
		panic("无法获取add函数")
	}
	// 调用函数
	result := C.int(C.call_function(addFuncPtr, C.int(1), C.int(2)))
	fmt.Println("Result:", result)
	fmt.Printf("1 + 2 = %d\n", result)

	// 转换函数类型
	// fn := func() uint32 {
	// 	ret, _, _ := syscall.SyscallN(uintptr(proc), 1, 2)
	// 	return uint32(ret)
	// }
	// // 调用函数
	// tickCount := fn()
	// fmt.Println("add:", tickCount)
}

//import "C"

// func main() {
// 	var handle unsafe.Pointer
// 	var addfunc func(int, int) int
// 	var err *C.char
// 	// handle = C.dlopen(C.CString("/lib/libm-2.6.1.so"), C.RTLD_NOW) linux
// 	handle = C.dlopen(C.CString("./test.so"), C.RTLD_LAZY)
// 	if handle == nil {
// 		err = C.dlerror()
// 		fmt.Printf("open lib error: %s\n", C.GoString(err))
// 		return
// 	}
// 	addfunc = *(*(func(int, int) int))(C.dlsym(handle, C.CString("add")))
// 	//addfunc = C.dlsym(handle, C.CString("add")).(func(int, int) int)
// 	if addfunc == nil {
// 		err = C.dlerror()
// 		fmt.Printf("symbol add not found , error: %s\n", C.GoString(err))
// 		return
// 	}
// 	fmt.Println(addfunc(1, 2))
// 	C.dlclose(handle)
// }

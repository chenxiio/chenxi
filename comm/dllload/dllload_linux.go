//go:build linux
// +build linux

package dllload

/*
 #cgo LDFLAGS: -L . -ldl -lstdc++
 #cgo CFLAGS: -I ./
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
*/
import "C"
import "unsafe"

func Foo() {
	println("This is running on linux.")
}
func LoadLibrary(libname string) uintptr {
	cStr := C.CString(libname + ".so")
	defer C.free(unsafe.Pointer(cStr))
	handle := C.dlopen(cStr, C.RTLD_LAZY)
	if handle == nil {
		panic("无法打开动态库")
	}

	return uintptr(handle)
}

func FreeLibrary(handle uintptr) (err error) {

	//C.dlclose(unsafe.Pointer(handle))
	return nil
}

func GetProcAddress(module uintptr, procname string) uintptr {
	cStr := C.CString(procname)
	defer C.free(unsafe.Pointer(cStr))

	FuncPtr := C.dlsym(unsafe.Pointer(module), cStr)
	if FuncPtr == nil {
		panic("无法获取add函数")
	}
	return uintptr(FuncPtr)
}

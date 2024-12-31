//go:build windows
// +build windows

package dllload

import (
	"syscall"
)

func Foo() {
	println("This is running on Windows.")
}
func LoadLibrary(libname string) uintptr {
	//p, err := filepath.Abs(libname + ".dll")
	lib, err := syscall.LoadLibrary(libname + ".dll")
	if err != nil || lib == 0 {
		panic(err)
	}
	return uintptr(lib)
}

func FreeLibrary(handle uintptr) (err error) {
	return syscall.FreeLibrary(syscall.Handle(handle))
}

func GetProcAddress(module uintptr, procname string) uintptr {
	p, err := syscall.GetProcAddress(syscall.Handle(module), procname)
	if err != nil || p == 0 {
		panic(err)
	}
	return uintptr(p)
}

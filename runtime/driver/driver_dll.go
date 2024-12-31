package driver

/*
#include <stdio.h>
#include <stdlib.h>
int add(void* func, int a, int b) {
    int (*fn)(int, int) = (int (*)(int, int))func;
    return fn(a, b);
}
int drv_start(void* func,char* param) {
	int (*fn)(char*);
 	fn = (int (*)(char*))func;
	return fn(param);
}

int drv_stop(void* func, char* param) {
    int (*fn)(char*) = (int (*)(char*))func;
    return fn(param);
}

int drv_read_int(void* func, const char* param, int* value) {
    int (*fn)(const char*, int*) = (int (*)(const char*, int*))func;
    return fn(param, value);
}

int drv_read_double(void* func, const char* param, double* value) {
    int (*fn)(const char*, double*) = (int (*)(const char*, double*))func;
    return fn(param, value);
}

int drv_read_string(void* func, const char* param, char* data_buf, int buf_size) {
    int (*fn)(const char*, char*, int) = (int (*)(const char*, char*, int))func;
    return fn(param, data_buf, buf_size);
}

int drv_write_int(void* func, const char* param, int value) {
    int (*fn)(const char*, int) = (int (*)(const char*, int))func;
    return fn(param, value);
}

int drv_write_double(void* func, const char* param, double value) {
    int (*fn)(const char*, double) = (int (*)(const char*, double))func;
    return fn(param, value);
}

int drv_write_string(void* func, const char* param, char* data) {
    int (*fn)(const char*, const char*) = (int (*)(const char*, const char*))func;
    return fn(param, data);
}
*/
import "C"
import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"github.com/chenxiio/chenxi/comm/dllload"
)

type DriverDll struct {
	plib uintptr
	//add_p              uintptr
	drv_start_p        uintptr
	drv_stop_p         uintptr
	drv_read_int_p     uintptr
	drv_read_double_p  uintptr
	drv_read_string_p  uintptr
	drv_write_int_p    uintptr
	drv_write_double_p uintptr
	drv_write_string_p uintptr
}

func NewDriverDll(path string, parm string) *DriverDll {

	dr := DriverDll{}
	fmt.Printf("")
	dr.plib = dllload.LoadLibrary(path)

	//dr.add_p = dllload.GetProcAddress(dr.plib, "add")
	dr.drv_start_p = dllload.GetProcAddress(dr.plib, "drv_start")
	dr.drv_stop_p = dllload.GetProcAddress(dr.plib, "drv_stop")
	dr.drv_read_int_p = dllload.GetProcAddress(dr.plib, "drv_read_int")
	dr.drv_read_double_p = dllload.GetProcAddress(dr.plib, "drv_read_double")
	dr.drv_read_string_p = dllload.GetProcAddress(dr.plib, "drv_read_string")
	dr.drv_write_int_p = dllload.GetProcAddress(dr.plib, "drv_write_int")
	dr.drv_write_double_p = dllload.GetProcAddress(dr.plib, "drv_write_double")
	dr.drv_write_string_p = dllload.GetProcAddress(dr.plib, "drv_write_string")

	if err := dr.Start(context.TODO(), parm); err != nil {
		panic(err)
	}
	return &dr
}

// func (d *DriverDll) Add(a, b int) int {
// 	return int((C.add(unsafe.Pointer(d.add_p), C.int(a), C.int(b))))
// }

func (d *DriverDll) Start(ctx context.Context, parm string) error {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))

	ret := int((C.drv_start(unsafe.Pointer(d.drv_start_p), pstr)))
	if ret != 0 {
		return fmt.Errorf("Start err ret = %d", ret)
	}
	return nil
}

func (d *DriverDll) Stop(ctx context.Context, parm string) error {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))

	ret := int((C.drv_stop(unsafe.Pointer(d.drv_stop_p), pstr)))
	if ret != 0 {
		return fmt.Errorf("Stop err ret = %d", ret)
	}
	return nil
}
func (d *DriverDll) ReadInt(ctx context.Context, parm string) (int32, error) {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	vint := C.int(0)
	ret := int((C.drv_read_int(unsafe.Pointer(d.drv_read_int_p), pstr, &vint)))
	if ret != 0 {
		return 0, fmt.Errorf("ReadInt err ret = %d", ret)
	}
	return int32(vint), nil
}
func (d *DriverDll) ReadString(ctx context.Context, parm string) (string, error) {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	buf := make([]byte, 256)
	ret := int(C.drv_read_string(unsafe.Pointer(d.drv_read_string_p), pstr, (*C.char)(unsafe.Pointer(&buf[0])), 256))
	if ret != 0 {
		return "", fmt.Errorf("ReadString err ret = %d", ret)
	}
	return strings.Trim(string(buf), "\x00"), nil
}
func (d *DriverDll) ReadDouble(ctx context.Context, parm string) (float64, error) {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	vdouble := C.double(0)
	ret := int(C.drv_read_double(unsafe.Pointer(d.drv_read_double_p), pstr, &vdouble))
	if ret != 0 {
		return 0, fmt.Errorf("ReadDouble err ret = %d", ret)
	}
	return float64(vdouble), nil
}

func (d *DriverDll) WriteInt(ctx context.Context, parm string, value int32) error {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	ret := int(C.drv_write_int(unsafe.Pointer(d.drv_write_int_p), pstr, C.int(value)))
	if ret != 0 {
		return fmt.Errorf("WriteInt err ret = %d", ret)
	}
	return nil
}
func (d *DriverDll) WriteString(ctx context.Context, parm string, value string) error {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	pdata := C.CString(value)
	defer C.free(unsafe.Pointer(pdata))
	ret := int(C.drv_write_string(unsafe.Pointer(d.drv_write_string_p), pstr, pdata))
	if ret != 0 {
		return fmt.Errorf("WriteString err ret = %d", ret)
	}
	return nil
}
func (d *DriverDll) WriteDouble(ctx context.Context, parm string, value float64) error {
	pstr := C.CString(parm)
	defer C.free(unsafe.Pointer(pstr))
	ret := int(C.drv_write_double(unsafe.Pointer(d.drv_write_double_p), pstr, C.double(value)))
	if ret != 0 {
		return fmt.Errorf("WriteDouble err ret = %d", ret)
	}
	return nil
}

// func (d *DriverDll) WriteBool(ctx context.Context, parm string, value bool) (bool, error) {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	var vint C.int
// 	if value {
// 		vint = 1
// 	}
// 	ret := int(C.drv_write_int(unsafe.Pointer(d.drv_write_int_p), pstr, vint))
// 	if ret != 0 {
// 		return false, fmt.Errorf("WriteBool err ret = %d", ret)
// 	}
// 	return value, nil
// }
// func (d *DriverDll) ReadBool(ctx context.Context, parm string) (bool, error) {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	vint := C.int(0)
// 	ret := int(C.drv_read_int(unsafe.Pointer(d.drv_read_int_p), pstr, &vint))
// 	if ret != 0 {
// 		return false, fmt.Errorf("ReadBool err ret = %d", ret)
// 	}
// 	if vint != 0 {
// 		return true, nil
// 	}
// 	return false, nil
// }

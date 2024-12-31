package driver

import (
	"context"
	"fmt"
	"plugin"

	"github.com/chenxiio/chenxi/api"
	"golang.org/x/exp/slog"
)

// type DriverPlugin struct {
// }

func NewDriverPlugin(path string, parm string) api.Drvapi {

	//dr := DriverPlugin{}

	p, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}
	// 声明一个公开实例变量，找到后直接使用
	d, err := p.Lookup("Driver")
	if err != nil {
		panic(err)
	}

	// 转换为Driver类型
	driver, ok := d.(api.Drvapi)
	if !ok {
		panic(fmt.Errorf("%s is not Drvapi", path))
	}

	// // 查找 Start 函数
	// fstart, err := p.Lookup("Start")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.Start = fstart.(func(ctx context.Context, parm string) error)

	// // 查找 Stop 函数
	// fstop, err := p.Lookup("Stop")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.Stop = fstop.(func(ctx context.Context, parm string) error)

	// // 查找 ReadInt 函数
	// freadint, err := p.Lookup("ReadInt")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.ReadInt = freadint.(func(ctx context.Context, parm string) (int, error))

	// // 查找 ReadString 函数
	// freadstr, err := p.Lookup("ReadString")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.ReadString = freadstr.(func(ctx context.Context, parm string) (string, error))

	// // 查找 ReadDouble 函数
	// freaddbl, err := p.Lookup("ReadDouble")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.ReadDouble = freaddbl.(func(ctx context.Context, parm string) (float64, error))

	// // 查找 WriteInt 函数
	// fwriteint, err := p.Lookup("WriteInt")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.WriteInt = fwriteint.(func(ctx context.Context, parm string, value int) error)

	// // 查找 WriteString 函数
	// fwritestr, err := p.Lookup("WriteString")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.WriteString = fwritestr.(func(ctx context.Context, parm string, value string) error)

	// // 查找 WriteDouble 函数
	// fwritedbl, err := p.Lookup("WriteDouble")
	// if err != nil {
	// 	panic(err)
	// }
	// dr.WriteDouble = fwritedbl.(func(ctx context.Context, parm string, value float64) error)

	if err := driver.Start(context.TODO(), parm); err != nil {
		slog.Error(err.Error())
	}
	return driver
}

// func (d *DriverDll) Add(a, b int) int {
// 	return int((C.add(unsafe.Pointer(d.add_p), C.int(a), C.int(b))))
// }

// func (d *DriverDll) Start(ctx context.Context, parm string) error {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))

// 	ret := int((C.drv_start(unsafe.Pointer(d.drv_start_p), pstr)))
// 	if ret != 0 {
// 		return fmt.Errorf("Start err ret = %d", ret)
// 	}
// 	return nil
// }

// func (d *DriverDll) Stop(ctx context.Context, parm string) error {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))

// 	ret := int((C.drv_stop(unsafe.Pointer(d.drv_stop_p), pstr)))
// 	if ret != 0 {
// 		return fmt.Errorf("Stop err ret = %d", ret)
// 	}
// 	return nil
// }
// func (d *DriverDll) ReadInt(ctx context.Context, parm string) (int, error) {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	vint := C.int(0)
// 	ret := int((C.drv_read_int(unsafe.Pointer(d.drv_read_int_p), pstr, &vint)))
// 	if ret != 0 {
// 		return 0, fmt.Errorf("ReadInt err ret = %d", ret)
// 	}
// 	return int(vint), nil
// }
// func (d *DriverDll) ReadString(ctx context.Context, parm string) (string, error) {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	buf := make([]byte, 256)
// 	ret := int(C.drv_read_string(unsafe.Pointer(d.drv_read_string_p), pstr, (*C.char)(unsafe.Pointer(&buf[0])), 256))
// 	if ret != 0 {
// 		return "", fmt.Errorf("ReadString err ret = %d", ret)
// 	}
// 	return strings.Trim(string(buf), "\x00"), nil
// }
// func (d *DriverDll) ReadDouble(ctx context.Context, parm string) (float64, error) {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	vdouble := C.double(0)
// 	ret := int(C.drv_read_double(unsafe.Pointer(d.drv_read_double_p), pstr, &vdouble))
// 	if ret != 0 {
// 		return 0, fmt.Errorf("ReadDouble err ret = %d", ret)
// 	}
// 	return float64(vdouble), nil
// }

// func (d *DriverDll) WriteInt(ctx context.Context, parm string, value int) error {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	ret := int(C.drv_write_int(unsafe.Pointer(d.drv_write_int_p), pstr, C.int(value)))
// 	if ret != 0 {
// 		return fmt.Errorf("WriteInt err ret = %d", ret)
// 	}
// 	return nil
// }
// func (d *DriverDll) WriteString(ctx context.Context, parm string, value string) error {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	pdata := C.CString(vaplugin
// 	defer C.free(unsafe.Pointer(pdata))
// 	ret := int(C.drv_write_string(unsafe.Pointer(d.drv_write_string_p), pstr, pdata))
// 	if ret != 0 {
// 		return fmt.Errorf("WriteString err ret = %d", ret)
// 	}
// 	return nil
// }
// func (d *DriverDll) WriteDouble(ctx context.Context, parm string, value float64) error {
// 	pstr := C.CString(parm)
// 	defer C.free(unsafe.Pointer(pstr))
// 	ret := int(C.drv_write_double(unsafe.Pointer(d.drv_write_double_p), pstr, C.double(value)))
// 	if ret != 0 {
// 		return fmt.Errorf("WriteDouble err ret = %d", ret)
// 	}
// 	return nil
// }

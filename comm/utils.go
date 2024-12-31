package comm

import (
	"fmt"
	"strconv"
	"strings"
)

func IODataConvert(dt, v string) (any, error) {
	if v == "" {
		return nil, nil
	}
	switch dt {
	case "int":
		return strconv.Atoi(v)

	case "string":
		// 将 i.Min 转换为 string 类型
		return v, nil
	case "double":
		// 将 i.Min 转换为 float64 类型
		return strconv.ParseFloat(v, 64)

	case "bool":
		// 将 i.Min 转换为 bool 类型
		return strconv.ParseBool(v)

	default:
		return nil, fmt.Errorf("unsupported data type: %s", dt)
	}

}

func UnmarshalUSlot(uslot string) (string, int, error) {

	us := strings.Split(uslot, ".")
	if len(us) < 2 {
		return "", -1, fmt.Errorf("格式错误或者wafer不在设备单元中： env.wid.%s", uslot)
	}
	slot, err := strconv.Atoi(us[len(us)-1])
	return us[0], slot, err
}

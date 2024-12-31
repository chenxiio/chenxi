package main

import (
	"encoding/json"
	"fmt"
)

type Menu int

const (
	MenuItem1 Menu = 1
	MenuItem2 Menu = 2
	MenuItem3 Menu = 3
)

func (m Menu) String() string {
	switch m {
	case MenuItem1:
		return "MenuItem1"
	case MenuItem2:
		return "MenuItem2"
	case MenuItem3:
		return "MenuItem3"
	default:
		return "Unknown"
	}
}

func main() {
	// 创建一个menu变量
	menu := MenuItem2

	// 将menu变量序列化为JSON字符串
	jsonStr, err := json.Marshal(menu)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON string:", string(jsonStr))
}

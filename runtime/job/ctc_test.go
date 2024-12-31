package job

import (
	"fmt"
	"testing"

	"github.com/chenxiio/chenxi/cfg"
)

func TestGetCTCInstance(t *testing.T) {
	cfg := &cfg.CTCCfg{} // 创建一个测试用的配置对象
	instance1, err := GetCTCInstance(cfg)
	if err != nil {
		t.Errorf("Failed to create CTC instance: %v", err)
	}
	instance2, err := GetCTCInstance(cfg)
	if err != nil {
		t.Errorf("Failed to create CTC instance: %v", err)
	}
	if instance1 != instance2 {
		t.Errorf("Expected instance1 and instance2 to be the same, but they are different")
	}
}
func TestDeeppath(t *testing.T) {
	mappath := map[string]string{
		"A": "B,C",
		"B": "D,E,F",
		"C": "F",
		"D": "G",
		"E": "H",
		"F": "I,J",
		"G": "K",
		"H": "L",
		"I": "M",
		"J": "N",
		"K": "O",
		"L": "P",
		"M": "Q",
		"N": "R",
		"O": "S",
		"P": "T",
		"Q": "U",
		"R": "V",
		"S": "W",
		"T": "X",
		"U": "Y",
		"Y": "Z",
		"W": "Z",
	}

	from := "A"
	to := "Z"
	rpath := []string{from}

	result := deeppath(mappath, rpath, from, to)
	fmt.Println(result)
}

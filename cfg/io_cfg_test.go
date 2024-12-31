package cfg

import (
	"encoding/xml"
	"os"
	"testing"
)

func TestReadIOConfigFile(t *testing.T) {
	//读取文件并检查内容
	file, err := os.Open("./temp/io_cfg.xml")
	if err != nil {
		t.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)
	var result IOCfg
	err = decoder.Decode(&result)
	if err != nil {
		t.Errorf("failed to decode config file: %v", err)
	}
	result.SaveConfigFile()
}

func TestSaveIOConfigFile(t *testing.T) {
	// 创建测试数据
	userMap := make(IODefines)
	userMap["IO2"] = IO{
		Name:  "IO2",
		DT:    "double",
		Cat:   "output",
		Dvid:  2,
		Svid:  2,
		Ecid:  2,
		Dfval: "0.0",
		Enum:  "",
		Unit:  "",
		Expr:  "",
		Min:   "",
		Max:   "",
		Pst:   "",
		Drv:   0,
		Pr:    "",
		Pw:    "",
		Rs:    200,
		Desc:  "",
	}
	userMap["IO1"] = IO{
		Name:  "IO1",
		DT:    "int",
		Cat:   "input",
		Dvid:  1,
		Svid:  1,
		Ecid:  1,
		Dfval: "0",
		Enum:  "",
		Unit:  "",
		Expr:  "",
		Min:   "",
		Max:   "",
		Pst:   "",
		Drv:   0,
		Pr:    "",
		Pw:    "",
		Rs:    200,
		Desc:  "",
	}

	cfg := IOCfg{
		// Help: IO{Name: "unit.test.test1",
		// DT:    "int,string,double,bool",
		// Cat:   "Memory,IO",
		// Dvid:  1,
		// Svid:  1,
		// Ecid:  1,
		// Dfval: "默认值",
		// Enum:  "",
		// Unit:  "",
		// Expr:  "",
		// Min:   "",
		// Max:   "",
		// Pst:   "",
		// Drv:   0,
		// Pr:    "读取时参数，例如",
		// Pw:    "写入时参数，例如",
		// Desc:  ""},
		Items: userMap,
		path:  "./temp/io_cfg.xml",
	}
	// 保存到文件
	err := cfg.SaveConfigFile()
	if err != nil {
		t.Errorf("failed to save config file: %v", err)
	}

	// if !reflect.DeepEqual(result, defines) {
	// 	t.Errorf("expected %v, got %v", defines, result)
	// }
}

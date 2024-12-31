package cfg

import (
	"encoding/xml"
	"os"
	"testing"
)

func Test_Module(t *testing.T) {
	// 定义测试数据
	xmlData := `
        <MODULE name="CM1" disp_name="CM1" sim_exec="CM.LP.exe"  exec="CM.LP.exe"  type="CM" alarm_index="58" priority="5" slot_count="25" cmd_svc="127.0.0.1:9026" />
    `

	// 测试XML转换为Golang结构体
	var module Module
	err := xml.Unmarshal([]byte(xmlData), &module)
	if err != nil {
		t.Errorf("Failed to unmarshal XML: %v", err)
	}

	// 测试Golang结构体转换为XML
	xmlData2, err := xml.Marshal(&module)
	if err != nil {
		t.Errorf("Failed to marshal XML: %v", err)
	}
	// if string(xmlData2) != xmlData {
	// 	t.Errorf("Failed to marshal XML correctly: %s", string(xmlData2))
	// }

	// 将XML输出到文件
	err = os.WriteFile("./temp/module_cfg.xml", xmlData2, 0644)
	if err != nil {
		t.Errorf("Failed to write XML to file: %v", err)
	}
}
func TestReadConfigFile(t *testing.T) {
	// 创建一个 Modules 实例
	modules := Modules{path: "./temp/module_cfg.xml"}
	// 创建一个测试用的配置文件
	// fileContent := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><modules><module name=\"module1\" disp_name=\"模块1\" sim_exec=\"/path/to/sim_exec\" exec=\"/path/to/exec\" type=\"type1\" alarm_index=\"1\" priority=\"2\" slot_count=\"3\" cmd_svc=\"cmd_svc1\"/><module name=\"module2\" disp_name=\"模块2\" sim_exec=\"/path/to/sim_exec\" exec=\"/path/to/exec\" type=\"type2\" alarm_index=\"2\" priority=\"1\" slot_count=\"4\" cmd_svc=\"cmd_svc2\"/></modules>"
	// err := os.WriteFile("config.xml", []byte(fileContent), 0644)
	// if err != nil {
	// 	t.Errorf("Create file failed: %v", err)
	// }
	// 调用 ReadConfigFile 方法读取配置文件
	err := modules.ReadConfigFile()
	if err != nil {
		t.Errorf("ReadConfigFile failed: %v", err)
	}
	// 检查读取结果是否正确
	if len(modules.Items) != 2 {
		t.Errorf("ReadConfigFile failed: expected 2 items, actual %d", len(modules.Items))
	}
	if modules.Items[0].Name != "module1" {
		t.Errorf("ReadConfigFile failed: expected name %s, actual %s", "module1", modules.Items[0].Name)
	}

	// // 删除测试用的配置文件
	// err = os.Remove("config.xml")
	// if err != nil {
	// 	t.Errorf("Remove file failed: %v", err)
	// }
}
func TestModuleSaveConfigFile(t *testing.T) {
	// 创建测试数据
	modules := Modules{
		// Help: Module{
		// 	ID:   1,
		// 	Name: "test_module_1",
		// 	API:  "driver,pm,cm,custom",
		// 	//LoadType: "dll ,plugin ,class,exe",
		// 	Path: "对象的路径",
		// 	Parm: "Start 的参数",
		// },
		Processes: Processs{"ioserver": Process{ProcessName: "ioserver", Url: "localhost:10600"}},
		Items: []Module{
			{
				ID:   1,
				Name: "test_module_1",
				API:  "test_type_1",
				//LoadType: "dll",
				Path: "test_path_1",
				Parm: "test_parm_1",
			},
			{
				ID:   2,
				Name: "test_module_2",
				API:  "test_type_2",
				//LoadType: "plugin",
				Path: "test_path_2",
				Parm: "test_parm_2",
			},
		},
		path: "./temp/module_cfg.xml",
	}
	// 保存到文件
	err := modules.SaveConfigFile()
	if err != nil {
		t.Errorf("failed to save config file: %v", err)
	}
	// // 读取文件并检查内容
	// file, err := os.Open("test.xml")
	// if err != nil {
	// 	t.Errorf("failed to open config file: %v", err)
	// }
	// defer file.Close()
	// decoder := xml.NewDecoder(file)
	// var result Modules
	// err = decoder.Decode(&result)
	// if err != nil {
	// 	t.Errorf("failed to decode config file: %v", err)
	// }
	// if !reflect.DeepEqual(result, modules) {
	// 	t.Errorf("expected %v, got %v", modules, result)
	// }
}

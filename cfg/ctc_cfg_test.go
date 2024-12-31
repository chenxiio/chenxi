package cfg

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

func TestSaveFileTm(t *testing.T) {
	tmConfigs := TMcfgs{
		"tm1": TMcfg{
			Name: "tm1",
			Type: 0,
			Arms: []string{"tm1.s1", "tm1.s2"},
		},
		"tm2": TMcfg{
			Name: "tm2",
			Type: 0,
			Arms: []string{"tm2.s1", "tm2.s2"},
		},
	}

	xmlData, err := xml.MarshalIndent(tmConfigs, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal XML:", err)
		return
	}
	fmt.Println(string(xmlData))

}
func TestMain(t *testing.T) {
	// 创建一个示例对象
	cfg := Carrier{
		Name: "CM1",

		Slots: []Slot{
			{
				Name: 1,
				//Priority: 1,
			},
			{
				Name: 2,
				//Priority: 2,
			},
		},
	}
	cfg2 := Carrier{
		Name: "CM2",

		Slots: []Slot{
			{
				Name: 1,
				//Priority: 1,
			},
			{
				Name: 2,
				//Priority: 2,
			},
		},
	}
	cfgs := CMcfgs{"CM1": cfg, "CM2": cfg2}

	// 序列化为XML
	xmlData, err := xml.MarshalIndent(cfgs, "", "  ")
	if err != nil {
		fmt.Println("XML序列化失败:", err)
		return
	}

	// 打印XML数据
	fmt.Println(string(xmlData))

	// 将XML数据写入文件
	file, err := os.Create("config.xml")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(xmlData)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}

	fmt.Println("XML序列化成功并写入文件")
}

func TestSaveFilectc(t *testing.T) {
	cfg := &CTCCfg{
		Group:         map[string][]string{"tmtest1": {"pmtest1", "cmtest1"}},
		Disable_paths: []string{"pmtest2-tmtest1"},
		Return_paths:  map[string]string{"pmtest1": "cmtest1", "path6": "p2"},
		Priority:      map[string]int{"tmtest1-pmtest1": 2, "tmtest1-cmtest1": 1},
		Interlocking: []*Interlock{&Interlock{Units: []string{"pmtest1", "pmtest2", "pmtest3"}, InParallel: 2, OutParallel: 3},
			&Interlock{Units: []string{"pmtest1-pmtest2-pmtest3-pmtest4-pmtest5-cmtest6"}}},
		path: "./temp/ctc.json",
	}
	err := cfg.SaveFile()
	if err != nil {
		t.Errorf("SaveFile() returned an error: %v", err)
	}
	// Read the saved file and compare the contents
	data, err := os.ReadFile(cfg.path)
	if err != nil {
		t.Errorf("Failed to read the saved file: %v", err)
	}
	var savedCfg CTCCfg
	err = json.Unmarshal(data, &savedCfg)
	if err != nil {
		t.Errorf("Failed to unmarshal the saved file: %v", err)
	}
	// Compare the original and saved configurations
	if len(cfg.Enable_paths) != len(savedCfg.Enable_paths) ||
		len(cfg.Disable_paths) != len(savedCfg.Disable_paths) ||
		len(cfg.Return_paths) != len(savedCfg.Return_paths) ||
		len(cfg.Priority) != len(savedCfg.Priority) {
		t.Errorf("Saved configuration does not match the original")
	}
	// Additional comparisons can be done here based on the requirements of your application
}

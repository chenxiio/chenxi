package cfg

import (
	"encoding/xml"
	"os"
	"testing"
)

func TestProcessSaveFile(t *testing.T) {
	// 创建一个 Procce 实例并设置测试数据
	procce := ProcessRecipe{
		Steps: []PStep{
			{Unit: []string{"TRS1|TRSTest", "TRS2|TRSTest"}, SubProcess: "emptprocess"},
			{Unit: []string{"PRS1|TRSTest", "PRS2|TRSTest"}},
		},
		path: "./temp/ProcessRecipe.xml",
	}
	// 调用 SaveFile 方法保存数据到文件
	err := procce.SaveFile()
	if err != nil {
		t.Errorf("SaveFile returned an error: %v", err)
	}
	// 读取保存的文件
	data, err := os.ReadFile(procce.path)
	if err != nil {
		t.Errorf("Failed to read saved file: %v", err)
	}
	// 解析 XML 数据
	var savedProcce ProcessRecipe
	err = xml.Unmarshal(data, &savedProcce)
	if err != nil {
		t.Errorf("Failed to unmarshal XML data: %v", err)
	}
	// 检查保存的数据是否与原始数据一致
	if len(savedProcce.Steps) != len(procce.Steps) {
		t.Errorf("Saved data has different number of steps")
	}
	for i := range procce.Steps {
		if len(savedProcce.Steps[i].Unit) != len(procce.Steps[i].Unit) {
			t.Errorf("Saved data has different number of sub-steps in step %d", i)
		}
		for j := range procce.Steps[i].Unit {
			if savedProcce.Steps[i].Unit[j] != procce.Steps[i].Unit[j] {
				t.Errorf("Saved data is different from original data in step %d, sub-step %d", i, j)
			}
		}
	}
	// 清理测试生成的文件
	err = os.Remove(procce.path)
	if err != nil {
		t.Errorf("Failed to remove test file: %v", err)
	}
}

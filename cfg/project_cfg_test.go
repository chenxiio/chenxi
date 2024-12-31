package cfg

import (
	"testing"
)

func TestProjecSaveConfigFile(t *testing.T) {
	// 创建测试数据
	project := Project{

		path: "./temp/project_cfg.xml",
		Alarms: KvItems{
			"alarm1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		DevicesSecurity: KvItems{
			"device1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		General: KvItems{
			"general1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		Notifications: KvItems{
			"notification1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		Reports: KvItems{
			"report1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		Scripts: KvItems{
			"script1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		Texts: KvItems{
			"text1": {
				Name:  "name1",
				Value: "value1",
			},
		},
		Views: KvItems{
			"view1": {
				Name:  "name1",
				Value: "value1",
			},
		},
	}
	// 保存到文件
	err := project.SaveConfigFile()
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
	// var result Project
	// err = decoder.Decode(&result)
	// if err != nil {
	// 	t.Errorf("failed to decode config file: %v", err)
	// }
	// if !reflect.DeepEqual(result, project) {
	// 	t.Errorf("expected %v, got %v", project, result)
	// }
}

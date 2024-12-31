package cfg

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"
)

func TestSaveFile(t *testing.T) {
	// 创建一个 UnitRecipe 实例
	recipe := &UnitRecipe{
		Type: "recipe",
		Info: Info{
			Items: []Item{
				{Name: "item1", Value: "value1"},
				{Name: "item2", Value: "value2"},
			},
		},
		Header: Header{
			Items: []Item{
				{Name: "header1", Value: "headerValue1"},
				{Name: "header2", Value: "headerValue2"},
			},
		},
		Steps: Steps{
			Step: Step{
				ID: 1,
				Items: []Item{
					{Name: "stepItem1", Value: "stepValue1"},
					{Name: "stepItem2", Value: "stepValue2"},
				},
			},
		},
		path: "./temp/unitrecipe.xml",
	}

	// 保存文件
	err := recipe.SaveFile()
	if err != nil {
		t.Errorf("保存文件时发生错误：%v", err)
	}

	// 读取保存的文件
	data, err := os.ReadFile(recipe.path)
	if err != nil {
		t.Errorf("读取保存的文件时发生错误：%v", err)
	}

	// 解析 XML 数据
	var savedRecipe UnitRecipe
	err = xml.Unmarshal(data, &savedRecipe)
	if err != nil {
		t.Errorf("解析 XML 数据时发生错误：%v", err)
	}

	// 比较原始 UnitRecipe 实例和保存后的实例是否相等
	if !reflect.DeepEqual(recipe, &savedRecipe) {
		t.Errorf("保存的 UnitRecipe 实例与原始实例不相等")
	}
}

package cfg

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	// 创建 UnitRecipeTPL 对象
	unitRecipe := UnitRecipeTPL{
		Nodes: map[string]Node{},
		path:  "data.xml",
	}
	unitRecipe.Nodes["HP"] = Node{
		Type: "HP",
		Header: Items{[]Item{
			{
				Disname:      "h Name 1",
				Name:         "Name 1",
				DefaultValue: "Default Value 1",
				Min:          "Min 1",
				Max:          "Max 1",
				Type:         "Type 1",
				Enum:         "Enum 1",
			},
			{
				Disname:      "h Name 2",
				Name:         "Name 2",
				DefaultValue: "Default Value 2",
				Min:          "Min 2",
				Max:          "Max 2",
				Type:         "Type 2",
			},
		}},
		Step: Items{[]Item{
			{
				Disname:      "Display Name 1",
				Name:         "Name 1",
				DefaultValue: "Default Value 1",
				Min:          "Min 1",
				Max:          "Max 1",
				Type:         "Type 1",
				Enum:         "Enum 1",
			},
			{
				Disname:      "Display Name 2",
				Name:         "Name 2",
				DefaultValue: "Default Value 2",
				Min:          "Min 2",
				Max:          "Max 2",
				Type:         "Type 2",
			},
			{
				Disname:      "Display Name 3",
				Name:         "Name 3",
				DefaultValue: "Default Value 3",
				Min:          "Min 3",
				Max:          "Max 3",
				Type:         "Type 3",
				Enum:         "Enum 3",
			}},
		},
	}
	unitRecipe.Nodes["CP"] = Node{
		Type: "CP",
		Header: Items{[]Item{
			{
				Disname:      "h Name 1",
				Name:         "Name 1",
				DefaultValue: "Default Value 1",
				Min:          "Min 1",
				Max:          "Max 1",
				Type:         "Type 1",
				Enum:         "Enum 1",
			},
			{
				Disname:      "h Name 2",
				Name:         "Name 2",
				DefaultValue: "Default Value 2",
				Min:          "Min 2",
				Max:          "Max 2",
				Type:         "Type 2",
			},
		}},
		Step: Items{[]Item{
			{
				Disname:      "Display Name 1",
				Name:         "Name 1",
				DefaultValue: "Default Value 1",
				Min:          "Min 1",
				Max:          "Max 1",
				Type:         "Type 1",
				Enum:         "Enum 1",
			},
			{
				Disname:      "Display Name 2",
				Name:         "Name 2",
				DefaultValue: "Default Value 2",
				Min:          "Min 2",
				Max:          "Max 2",
				Type:         "Type 2",
			},
			{
				Disname:      "Display Name 3",
				Name:         "Name 3",
				DefaultValue: "Default Value 3",
				Min:          "Min 3",
				Max:          "Max 3",
				Type:         "Type 3",
				Enum:         "Enum 3",
			}},
		},
	}
	// 测试 SaveFile 方法
	err := unitRecipe.SaveFile()
	if err != nil {
		fmt.Println("保存文件时出错：", err)
		return
	}

	fmt.Println("文件保存成功！")
	// // 读取文件
	// unitRecipe, err := ReadFile("data.xml")
	// if err != nil {
	// 	fmt.Println("读取文件时出错：", err)
	// 	return
	// }

	// // 修改 unitRecipe 对象

	// // 保存文件
	// err = SaveFile("output.xml", unitRecipe)
	// if err != nil {
	// 	fmt.Println("保存文件时出错：", err)
	// 	return
	// }

	// fmt.Println("文件保存成功！")
}

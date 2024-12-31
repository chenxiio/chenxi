package recipe

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetUnitRecipeList(t *testing.T) {
	// 创建一个临时文件夹作为测试数据
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时文件夹中创建一些测试文件
	testFiles := []string{"file1.txt", "file2.txt", "file3", "folder1", "folder2"}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if file == "folder1" || file == "folder2" {
			err := os.Mkdir(filePath, 0755)
			if err != nil {
				t.Fatalf("Failed to create test folder: %v", err)
			}
		} else {
			_, err := os.Create(filePath)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
	}

	// 创建 RecipeEdit 实例并设置测试数据
	r := Recipe{
		processdatapath: tempDir,
	}

	// 调用 GetUnitRecipeList 方法获取文件列表
	fileList, err := r.ReadProcessRecipeList(context.TODO())
	if err != nil {
		t.Fatalf("GetUnitRecipeList returned an error: %v", err)
	}

	// 检查返回的文件列表是否正确
	expectedList := []string{"file1", "file2", "file3"}
	if len(fileList) != len(expectedList) {
		t.Fatalf("GetUnitRecipeList returned incorrect number of files. Expected %d, got %d", len(expectedList), len(fileList))
	}

	for i, file := range expectedList {
		if fileList[i] != file {
			t.Errorf("GetUnitRecipeList returned incorrect file name at index %d. Expected %s, got %s", i, file, fileList[i])
		}
	}
}

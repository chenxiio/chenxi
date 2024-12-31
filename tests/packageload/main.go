package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

func main() {
	typeName := "mypackage.Person"
	pkgPath := "mypackage"
	// 解析包
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("解析包失败：", err)
		return
	}
	// 检查包
	cfg := types.Config{}
	for _, pkg := range pkgs {
		files := make([]*ast.File, 0, len(pkg.Files))
		for _, file := range pkg.Files {
			files = append(files, file)
		}
		info := &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		}
		_, err := cfg.Check(pkg.Name, fset, files, info)
		if err != nil {
			fmt.Println("检查包失败：", err)
			return
		}

		// 获取类型信息
		//pType := info.ObjectOf(pkg.Scope.Lookup(typeName).(*ast.Object).Name).Type()
		pType := pkg.Scope.Lookup(typeName).Type
		if pType == nil {
			fmt.Println("找不到类型：", typeName)
			return
		}
		fmt.Println("类型信息：", pType)
	}
}

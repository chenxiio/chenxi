package main

import (
	"fmt"
)

func main() {

	pkgPath := "importpkg/testpkg"
	pkg, err := importPackage(pkgPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Imported package %s: %v", pkgPath, pkg)
}
func importPackage(pkgPath string) (interface{}, error) {
	return importWithDepth(pkgPath, 0)
}
func importWithDepth(pkgPath string, depth int) (interface{}, error) {
	if depth > 10 {
		return nil, fmt.Errorf("maximum import depth exceeded")
	}
	pkg, err := importer(pkgPath).Import(pkgPath)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
func importer(pkgPath string) interface {
	Import(string) (interface{}, error)
} {
	return &importerImpl{}
}

type importerImpl struct{}

func (i *importerImpl) Import(pkgPath string) (interface{}, error) {
	return importWithDepth(pkgPath, 1)
}

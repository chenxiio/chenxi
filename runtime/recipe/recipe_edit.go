package recipe

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/chenxiio/chenxi/cfg"
)

type Recipe struct {
	//datapath        string
	unitdatapath    string
	processdatapath string
	tplpath         string
}

func NewRecipe(path string) *Recipe {
	return &Recipe{tplpath: path + "tpl/", unitdatapath: path + "data/unit/", processdatapath: path + "data/process/"}
}

func (r *Recipe) ReadUnitRecipeList(ctx context.Context, t string) ([]string, error) {
	files, err := ioutil.ReadDir(r.unitdatapath + t)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		}
	}

	return fileNames, nil
}
func (r *Recipe) ReadProcessRecipeList(ctx context.Context) ([]string, error) {
	files, err := ioutil.ReadDir(r.processdatapath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		}
	}

	return fileNames, nil
}
func (r *Recipe) ReadUnitRecipe(ctx context.Context, t, name string) (cfg.UnitRecipe, error) {
	rcp := cfg.NewUnitRecipe(r.unitdatapath + t + "/" + name)
	err := rcp.ReadFile()
	return rcp, err
}
func (r *Recipe) ReadProcessRecipe(ctx context.Context, name string) (cfg.ProcessRecipe, error) {
	rcp := cfg.NewProcessRecipe(r.processdatapath+name, r.unitdatapath)
	err := rcp.ReadFile()
	return rcp, err
}

func (r *Recipe) SaveUnitRecipe(ctx context.Context, t, name string, rcp cfg.UnitRecipe) error {
	rcp.Setpath(r.unitdatapath + t + "/" + name)
	return rcp.SaveFile()
}
func (r *Recipe) SaveProcessRecipe(ctx context.Context, name string, rcp cfg.ProcessRecipe) error {
	rcp.Setpath(r.processdatapath + name)
	return rcp.SaveFile()
}

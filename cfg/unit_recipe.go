package cfg

import (
	"encoding/xml"
	"os"
)

type UnitRecipe struct {
	Type   string `xml:"type,attr"`
	Info   Info   `xml:"Info"`
	Header Header `xml:"Header"`
	Steps  Steps  `xml:"Steps"`
	path   string
}
type Info struct {
	Items []Item `xml:"Item"`
}

type Header struct {
	Items []Item `xml:"Item"`
}

type Steps struct {
	Step Step `xml:"Step"`
}

type Step struct {
	ID    int    `xml:"id,attr"`
	Items []Item `xml:"Item"`
}

//	type Item struct {
//		Name  string `xml:"name,attr"`
//		Value string `xml:"value,attr"`
//	}
func NewUnitRecipe(path string) UnitRecipe {

	return UnitRecipe{path: path}
}
func (u *UnitRecipe) Setpath(p string) {
	u.path = p
}
func (u *UnitRecipe) SaveFile() error {
	data, err := xml.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(u.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (u *UnitRecipe) ReadFile() error {
	data, err := os.ReadFile(u.path)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, u)
	if err != nil {
		return err
	}
	return nil
}

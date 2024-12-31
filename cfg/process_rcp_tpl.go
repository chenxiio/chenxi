package cfg

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// type PStepItem struct {
// 	UName string
// 	Rcp   string
// }

type PStep struct {
	//Name       int      `xml:"name,attr"` // step name
	Unit       []string // [TRS|TRSTest,TRS2|TRSTest]
	Curmat     string   `xml:"-"` //
	IsSubStart bool     `xml:"-"`
	SubProcess string   `xml:"subprocess,attr"` // 子流程 ，常用于转移空foup，载具
}
type ProcessRecipe struct {
	Steps    []PStep `xml:"Step"`
	path     string
	unitpath string
	unitrcps map[string]UnitRecipe
}

func NewProcessRecipe(path string, unitpath string) ProcessRecipe {

	return ProcessRecipe{path: path, unitpath: unitpath, unitrcps: map[string]UnitRecipe{}}
}

func (u *ProcessRecipe) GetUnitRcp(step int, unitname string, unittype string) (UnitRecipe, error) {
	key := fmt.Sprintf("%d-%s", step, unitname)
	if urcp, ok := u.unitrcps[key]; ok {
		return urcp, nil
	}

	for k, v := range u.Steps {
		if k+1 == step {
			for _, v1 := range v.Unit {
				ur := strings.Split(v1, "|")
				if len(ur) == 2 {
					if ur[0] == unitname {
						urcp := NewUnitRecipe(u.unitpath + unittype + "/" + ur[1])
						err := urcp.ReadFile()
						u.unitrcps[key] = urcp
						return urcp, err
					}
				}
			}
		}
	}
	return UnitRecipe{}, fmt.Errorf("getUnitRcp(%d,%s) not fond", step, unitname)
}

// func (u *ProcessRecipe) GetRcpnameByunit(step int, unitname string) (string, error) {
// 	// if urcp, ok := u.unitrcps[unitname]; ok {
// 	// 	return urcp, nil
// 	// }
// 	unittype :=
// 	for _, v := range u.Steps {
// 		if v.Name == step {
// 			for _, v1 := range v.Unit {
// 				ur := strings.Split(v1, "|")
// 				if len(ur) == 2 {
// 					if ur[0] == unittype {
// 						// urcp := NewUnitRecipe(u.unitpath + unitname + "/" + ur[1])
// 						// err := urcp.ReadFile()

//							return ur[1], nil
//						}
//					}
//				}
//			}
//		}
//		return "", fmt.Errorf("GetUnitRcp(%s,%s) not fond", step, unittype)
//	}
func (u *ProcessRecipe) Setpath(p string) {
	u.path = p
}
func (u *ProcessRecipe) SaveFile() error {
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

func (u *ProcessRecipe) ReadFile() error {
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

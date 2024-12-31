package cfg

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type LoadType int

const (
	DLL LoadType = iota
	Plugin
	Class
)

type Process struct {
	ProcessName string `xml:"process_name,attr"`
	Url         string `xml:"url,attr"`
}
type Processs map[string]Process

func (m Processs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for _, v := range m {
		e.Encode(v)
	}

	return e.EncodeToken(start.End())
}

func (m *Processs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Processs{}
	for {
		var e Process

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.ProcessName] = e
	}
	return nil
}

// type Exes struct {
// }

// loadtype: dll ,plugin ,class,exe
type Module struct {
	ID          int    `xml:"id,attr"`
	ProcessName string `xml:"process_name,attr"`
	Name        string `xml:"name,attr"`
	API         string `xml:"api,attr"`
	Type        string `xml:"type,attr"`
	Path        string `xml:"path,attr"`
	Parm        string `xml:"parm,attr"`
	Slot_count  int    `xml:"slot_count,attr"`
	Priority    int    `xml:"priority,attr"`
	PositionNo  int    `xml:"psno,attr"`
	utype       string
}

// func (m Module) GetType() string {
// 	if m.utype == "" {
// 		at := strings.Split(m.API, ".")
// 		if len(at) == 2 {
// 			m.utype = at[1]
// 		}
// 	}
// 	return m.utype
// }

type Modules struct {
	// Help    Module   `xml:"Help"`
	Items     []Module `xml:"Module"`
	Processes Processs
	CMcfgs    CMcfgs
	PMcfgs    PMcfgs
	TMcfgs    TMcfgs
	path      string
}

func NewCfgModules(bsdir string) *Modules {
	return &Modules{path: fmt.Sprintf("%scfg/module_cfg.xml", bsdir)}
}

func (m *Modules) GetCfgByUnit(name string) (Module, error) {

	for _, v := range m.Items {
		if v.Name == name {
			return v, nil
		}
	}
	return Module{}, fmt.Errorf("GetCfgByUnit(%s) not fond", name)
}

// func (m *Modules) GetCMNamebyslot(slot string) (string, error) {

// 	for _, v := range m.CMcfgs {
// 		if v.Slot_name == slot {
// 			return v.Name, nil
// 		}
// 	}
// 	return "", fmt.Errorf("GetCMNamebyslot(%s) not fond", slot)
// }

//	func (m *Modules) GetSlotNamebycm(cm string) (string, error) {
//		if v, ok := m.CMcfgs[cm]; ok {
//			return v.Slot_name, nil
//		}
//		// for _, v := range m.CMcfgs {
//		// 	if v.Name == cm {
//		// 		return v.Slot_name, nil
//		// 	}
//		// }
//		return "", fmt.Errorf("GetSlotNamebycm(%s) not fond", cm)
//	}
func (m *Modules) SaveConfigFile() error {
	configFile, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	encoder := xml.NewEncoder(configFile)
	encoder.Indent("", " ")
	err = encoder.Encode(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) ReadConfigFile() error {
	configFile, err := os.Open(m.path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := xml.NewDecoder(configFile)
	err = decoder.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

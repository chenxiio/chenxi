package cfg

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/chenxiio/chenxi/comm"
)

// type CM3 struct {
// 	Name  string `json:"name"`
// 	DT    string `json:"dt"`
// 	Cat   string `json:"cat"`
// 	Dvid  string `json:"dvid"`
// 	Svid  string `json:"svid"`
// 	Ecid  string `json:"ecid"`
// 	Dfval string `json:"dfval"`
// 	Enum  string `json:"enum"`
// 	Unit  string `json:"unit"`
// 	Expr  string `json:"expr"`
// 	Min   string `json:"min"`
// 	Max   string `json:"max"`
// 	Pst   string `json:"pst"`
// 	Drv   string `json:"drv"`
// 	Pr    string `json:"pr"`
// 	Pw    string `json:"pw"`
// 	Desc  string `json:"desc"`
// }

type IO struct {
	Name  string `xml:"name,attr"`
	DT    string `xml:"dt,attr"`
	Cat   string `xml:"cat,attr"`
	Dvid  int    `xml:"dvid,attr"`
	Svid  int    `xml:"svid,attr"`
	Ecid  int    `xml:"ecid,attr"`
	Dfval string `xml:"dfval,attr"`
	Enum  string `xml:"enum,attr"`
	Unit  string `xml:"unit,attr"`
	Expr  string `xml:"expr,attr"`
	Min   string `xml:"min,attr"`
	Max   string `xml:"max,attr"`
	Pst   string `xml:"pst,attr"`
	Drv   int    `xml:"drv,attr"`
	Pr    string `xml:"pr,attr"`
	Pw    string `xml:"pw,attr"`
	Rs    int    `xml:"rs,attr"`
	Rcd   int    `xml:"rcd,attr"`
	Desc  string `xml:"desc,attr"`

	min   any
	max   any
	dfval any
}

type IODefines map[string]IO
type IOCfg struct {
	// Help  IO        `xml:"Help"`
	Items IODefines `xml:"Items"`
	path  string
}

func (i *IO) GetMin() (any, error) {
	//int,string,double,bool
	var err error
	if i.min == nil {
		i.min, err = comm.IODataConvert(i.DT, i.Min)
	}
	return i.min, err
}
func (i *IO) GetMax() (any, error) {
	var err error
	if i.max == nil {
		i.max, err = comm.IODataConvert(i.DT, i.Max)
	}
	return i.max, err
}
func (i *IO) GetDfval() (any, error) {
	var err error
	if i.dfval == nil {
		i.dfval, err = comm.IODataConvert(i.DT, i.Dfval)
	}
	return i.dfval, err
}

func (m IODefines) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *IODefines) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = IODefines{}
	for {
		var e IO

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.Name] = e
	}
	return nil
}
func (m *IOCfg) SaveConfigFile() error {
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

func (m *IOCfg) ReadConfigFile() error {
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

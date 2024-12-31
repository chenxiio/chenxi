package cfg

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/chenxiio/chenxi/comm"
)

type Nodes map[string]Node

type UnitRecipeTPL struct {
	Nodes Nodes
	path  string
}

type Items struct {
	Items []Item `xml:"Item"`
}
type Node struct {
	Type   string `xml:"type,attr"`
	Header Items
	Step   Items
}

type Item struct {
	Disname      string `xml:"disname,attr,omitempty"`
	Name         string `xml:"name,attr,omitempty"`
	DT           string `xml:"dt,attr,omitempty"`
	DefaultValue string `xml:"default_value,attr,omitempty"`
	Value        string `xml:"value,attr,omitempty"`
	Min          string `xml:"min,attr,omitempty"`
	Max          string `xml:"max,attr,omitempty"`
	Type         string `xml:"type,attr,omitempty"`
	Enum         string `xml:"enum,attr,omitempty"`

	min   any
	max   any
	dfval any
}

func (m Nodes) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *Nodes) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Nodes{}
	for {
		var e Node

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.Type] = e
	}
	return nil
}

func (i *Item) GetMin() (any, error) {
	//int,string,double,bool
	var err error
	if i.min == nil {
		i.min, err = comm.IODataConvert(i.DT, i.Min)
	}
	return i.min, err
}
func (i *Item) GetMax() (any, error) {
	var err error
	if i.max == nil {
		i.max, err = comm.IODataConvert(i.DT, i.Max)
	}
	return i.max, err
}
func (i *Item) GetDfval() (any, error) {
	var err error
	if i.dfval == nil {
		i.dfval, err = comm.IODataConvert(i.DT, i.DefaultValue)
	}
	return i.dfval, err
}

func (u *UnitRecipeTPL) SaveFile() error {
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

func (u *UnitRecipeTPL) ReadFile() error {
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

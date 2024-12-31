package stringmap

import (
	"encoding/xml"
	"io"
)

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
	Desc  string `xml:"desc,attr"`
	min   any
	max   any
	dfval any
}
type StringMap map[string]IO

// type xmlMapEntry struct {
// 	XMLName xml.Name
// 	Value   string `xml:",chardata"`
// }

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
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

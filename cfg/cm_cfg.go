package cfg

import (
	"encoding/xml"
	"io"
)

type Slot struct {
	Name int `xml:"name,attr"`
	//Priority int    `xml:"priority,attr"`
	WaferId string `xml:"-"`
	State   string `xml:"-"` // 正在搬运
}

// 由小到大排序
type Slots []Slot

// type Carrier struct {
// 	Slots Slots
// }

type Carrier struct {
	Name  string `xml:"name,attr"`
	Slots Slots  `xml:"slot"`
}

type CMcfgs map[string]Carrier

func (m CMcfgs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *CMcfgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = CMcfgs{}
	for {
		var e Carrier

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

// func (m Slots) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	if len(m) == 0 {
// 		return nil
// 	}

// 	err := e.EncodeToken(start)
// 	if err != nil {
// 		return err
// 	}

// 	for _, v := range m {
// 		e.Encode(v)
// 	}

// 	return e.EncodeToken(start.End())
// }

// func (m *Slots) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	*m = Slots{}
// 	for {
// 		var e Slot
// 		err := d.Decode(&e)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return err
// 		}
// 		(*m)[e.Name] = e
// 	}
// 	return nil
// }

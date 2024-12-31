package cfg

import (
	"encoding/xml"
	"io"
)

type PMcfg struct {
	Name string `xml:"name,attr"`
	// Slot_count int    `xml:"slot_count,attr,omitempty"`
}

type PMcfgs map[string]PMcfg

func (m PMcfgs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *PMcfgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = PMcfgs{}
	for {
		var e PMcfg

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

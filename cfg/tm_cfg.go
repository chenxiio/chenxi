package cfg

import (
	"encoding/xml"
	"io"
)

type TMType int

const (
	WAFER   TMType = 0
	CARRIER TMType = 1
)

// type : CARRIER， WAFER
type TMcfg struct {
	Name string   `xml:"name,attr"`
	Type TMType   `xml:"type,attr"`
	Arms []string `xml:"Arm"` // 1
}

type TMcfgs map[string]TMcfg

func (m TMcfgs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *TMcfgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = TMcfgs{}
	for {
		var e TMcfg

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

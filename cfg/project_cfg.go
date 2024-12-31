package cfg

import (
	"encoding/xml"
	"io"
	"os"
)

type Kv struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Project struct {
	Alarms          KvItems
	DevicesSecurity KvItems
	General         KvItems
	Notifications   KvItems
	Reports         KvItems
	Scripts         KvItems
	Texts           KvItems
	Views           KvItems

	path string
}

func (p *Project) SaveConfigFile() error {
	configFile, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	encoder := xml.NewEncoder(configFile)
	encoder.Indent("", " ")
	err = encoder.Encode(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) ReadConfigFile() error {
	configFile, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := xml.NewDecoder(configFile)
	err = decoder.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

type KvItems map[string]Kv

func (m KvItems) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *KvItems) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = KvItems{}
	for {
		var e Kv

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

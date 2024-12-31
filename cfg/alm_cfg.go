package cfg

import (
	"encoding/xml"
	"io"
	"os"
)

type AlmClearType byte

const (
	Abort  = AlmClearType(1)
	Ignore = AlmClearType(2)
	Clear  = AlmClearType(3)
	Retry  = AlmClearType(4)
	ALL    = AlmClearType(255)
)

type Alm struct {
	AlarmID     int    `xml:"AlarmID,attr"`
	Type        string `xml:"Type,attr"`
	Retry       int    `xml:"Retry,attr"`
	Abort       int    `xml:"Abort,attr"`
	Ignore      int    `xml:"Ignore,attr"`
	Clear       int    `xml:"Clear,attr"`
	Level       int    `xml:"Level,attr"`
	Flag        int    `xml:"Flag,attr"`
	Enable      int    `xml:"Enable,attr"`
	AlarmName   string `xml:"AlarmName,attr"`
	AlarmDesc   string `xml:"AlarmDesc,attr"`
	AlarmDetail string `xml:"AlarmDetail,attr"`
}

type AlmItems map[int]Alm

func (m AlmItems) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

func (m *AlmItems) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = AlmItems{}
	for {
		var e Alm

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.AlarmID] = e
	}
	return nil
}

type Alarms struct {
	Alarms AlmItems
	path   string
}

func (p *Alarms) SaveConfigFile() error {
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

func (p *Alarms) ReadConfigFile() error {
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

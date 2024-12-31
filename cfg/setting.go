package cfg

import (
	"encoding/xml"
	"os"
)

type Setting struct {
	Current_project string `xml:"current_project,attr"`

	path string
}

func (p *Setting) SaveConfigFile() error {
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

func (p *Setting) ReadConfigFile() error {
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

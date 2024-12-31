package main

import (
	"encoding/xml"
	"fmt"
)

type DRV struct {
	XMLName     xml.Name `xml:"DRV"`
	ID          string   `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Enable      string   `xml:"enable,attr"`
	File        string   `xml:"file,attr"`
	ParamStart  string   `xml:"param_start,attr"`
	ParamStop   string   `xml:"param_stop,attr"`
	Safe        string   `xml:"safe,attr"`
	Cache       string   `xml:"cache,attr"`
	Description string   `xml:"desc,attr"`
}

func main() {
	xmlString := `<DRV id="316" name="RFID_4" enable="False" file="CIDRW_V640_4.dll" param_start="9,38400,8,1,2,0,0,0,1,10000" param_stop="" safe="None" cache="None" desc="RFID1-loadport1-12" />`
	var drv DRV
	err := xml.Unmarshal([]byte(xmlString), &drv)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%+v\n", drv)
}

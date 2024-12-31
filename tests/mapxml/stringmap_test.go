package stringmap

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	userMap := make(map[string]IO)
	userMap["IO2"] = IO{
		Name:  "IO2",
		DT:    "double",
		Cat:   "output",
		Dvid:  2,
		Svid:  2,
		Ecid:  2,
		Dfval: "0.0",
		Enum:  "",
		Unit:  "",
		Expr:  "",
		Min:   "",
		Max:   "",
		Pst:   "",
		Drv:   0,
		Pr:    "",
		Pw:    "",
		Desc:  "",
	}
	userMap["IO1"] = IO{
		Name:  "IO1",
		DT:    "int",
		Cat:   "input",
		Dvid:  1,
		Svid:  1,
		Ecid:  1,
		Dfval: "0",
		Enum:  "",
		Unit:  "",
		Expr:  "",
		Min:   "",
		Max:   "",
		Pst:   "",
		Drv:   0,
		Pr:    "",
		Pw:    "",
		Desc:  "",
	}

	buf, _ := xml.MarshalIndent(StringMap(userMap), "", " ")
	fmt.Println(string(buf))

	stringMap := make(map[string]IO)
	err := xml.Unmarshal(buf, (*StringMap)(&stringMap))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stringMap)
}

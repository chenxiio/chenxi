package models

import (
	"bytes"
	"encoding/binary"
)

type His struct {
	Parm       string `json:"Parm"`
	Value      any    `json:"Value"`
	CreateTime int64  `json:"CreateTime"`
	//Cvalue     any    `json:"Value"`
}

// func (h *His) GetValue() any {
// 	if h.Cvalue == nil {
// 		h.Cvalue = ConvertValue(h.Value)
// 	}
// 	return h.Cvalue
// }

type IO_TYPE byte

const (
	Unknown = IO_TYPE(iota)
	//	IO_TYPE_BOOL   = IO_TYPE(1)
	IO_TYPE_INT    = IO_TYPE(2)
	IO_TYPE_DOUBLE = IO_TYPE(3)
	IO_TYPE_STRING = IO_TYPE(4)
)

func ConvertValue(existingValue []byte) interface{} {
	switch IO_TYPE(existingValue[0]) {
	// case IO_TYPE_BOOL:
	// 	var value bool
	// 	binary.Read(bytes.NewBuffer(existingValue[1:]), binary.BigEndian, &value)
	// 	return value
	case IO_TYPE_INT:
		var value int32
		binary.Read(bytes.NewBuffer(existingValue[1:]), binary.BigEndian, &value)
		return value
	case IO_TYPE_DOUBLE:
		var value float64
		binary.Read(bytes.NewBuffer(existingValue[1:]), binary.BigEndian, &value)
		return value
	case IO_TYPE_STRING:
		// var value string
		// bytes.NewBuffer(existingValue[1:]).ReadString()
		return string(existingValue[1:])
	default:
		return nil
	}
}

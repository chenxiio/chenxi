package tests

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestMain() {

}
func matchPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
func TestMatchPrefix(t *testing.T) {
	s := "Hello, world!"
	prefix := "Hello"
	if !matchPrefix(s, prefix) {
		t.Errorf("'%s' should start with '%s'", s, prefix)
	}
	s = "Golang is awesome!"
	prefix = "Python"
	if matchPrefix(s, prefix) {
		t.Errorf("'%s' should not start with '%s'", s, prefix)
	}
}

func Test_json(t *testing.T) {
	m := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", b)
}

func Test_any2bytes(t *testing.T) {
	// 将int转换为[]byte
	var buf bytes.Buffer
	i := 123
	err := binary.Write(&buf, binary.LittleEndian, i)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(buf.Bytes())

	// 将string转换为[]byte
	var buf2 bytes.Buffer
	s := "hello"
	err = binary.Write(&buf2, binary.LittleEndian, s)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(buf2.Bytes())

	// 将float32转换为[]byte
	var buf3 bytes.Buffer
	f := float32(3.14)
	err = binary.Write(&buf3, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(buf3.Bytes())

	// 将bool转换为[]byte
	var buf4 bytes.Buffer
	b1 := true
	err = binary.Write(&buf4, binary.LittleEndian, b1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(buf4.Bytes())
}

func TestTiker(t *testing.T) {

	fmt.Println("定时器已启动")
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("定时器触发")
		}
	}

}

package driver

import (
	"context"
	"testing"
)

var projpath = "../build/testproj/"

func TestAdd(t *testing.T) {
	d := NewDriverDll(projpath+"bin/drivers/testdriver", "")

	//fmt.Printf("Add:%d\n", d.Add(1, 2))

	// test ReadString
	expected := "hello world"
	if err := d.WriteString(context.Background(), "test_string", expected); err != nil {
		t.Errorf("WriteString err: %v", err)
	}
	actual, err := d.ReadString(context.Background(), "test_string")
	if err != nil {
		t.Errorf("ReadString err: %v", err)
	}
	if actual != expected {
		t.Errorf("ReadString failed, expected: %s, got: %s", expected, actual)
	}

	// test ReadDouble
	expectedD := 3.14
	if err := d.WriteDouble(context.Background(), "test_double", expectedD); err != nil {
		t.Errorf("WriteDouble err: %v", err)
	}
	actualD, err := d.ReadDouble(context.Background(), "test_double")
	if err != nil {
		t.Errorf("ReadDouble err: %v", err)
	}
	if actualD != expectedD {
		t.Errorf("ReadDouble failed, expected: %f, got: %f", expectedD, actualD)
	}

	// test WriteInt
	expectedI := int32(42)
	err = d.WriteInt(context.Background(), "test_int", expectedI)
	if err != nil {
		t.Errorf("WriteInt err: %v", err)
	}
	actualI, err := d.ReadInt(context.Background(), "test_int")
	if actualI != expectedI {
		t.Errorf("WriteInt failed, expected: %d, got: %d", expectedI, actualI)
	}

}

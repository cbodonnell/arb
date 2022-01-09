package arb

import (
	"os"
	"strings"
	"testing"
)

// TestRead tests reading from an io.Reader
func TestRead(t *testing.T) {
	jsonString := "{ \"key\": \"value\" }"
	expected := "value"
	reader := strings.NewReader(jsonString)
	arb, err := Read(reader)
	if err != nil {
		t.Fatal("failed to read json string from reader:", err)
	}
	if arb["key"] != expected {
		t.Fatalf("unexpected value for key %s, expected %s", arb["key"], expected)
	}
}

// TestReadBytes tests reading from a byte slice
func TestReadBytes(t *testing.T) {
	jsonString := "{ \"key\": \"value\" }"
	expected := "value"
	arb, err := ReadBytes([]byte(jsonString))
	if err != nil {
		t.Fatal("failed to read json string from bytes:", err)
	}
	if arb["key"] != expected {
		t.Fatalf("unexpected value for key %s, expected %s", arb["key"], expected)
	}
}

// TestWrite tests writing from an io.Write
func TestWrite(t *testing.T) {
	arb := New()
	arb["hello"] = "world"
	arb["number"] = 2222
	arb["nothing"] = nil
	child := New()
	child["isChild"] = true
	arb["child"] = child

	err := arb.Write(os.Stdout)
	if err != nil {
		t.Fatal("failed to write:", err)
	}
}

// TestReadBytes tests reading from a byte slice
func TestGetString(t *testing.T) {
	arb := New()
	arb["hello"] = "world"
	arb["number"] = 2222
	result, err := arb.GetString("hello")
	expected := "world"
	if err != nil {
		t.Fatal("failed get string:", err)
	}
	if result != expected {
		t.Fatalf("result was %s, expected %s", result, expected)
	}
	if _, err := arb.GetString("number"); err == nil {
		t.Fatalf("number %d is not a string", arb["number"])
	}
}

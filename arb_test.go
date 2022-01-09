package arb

import (
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
		t.Fatal("failed to read json string from reader", err)
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
		t.Fatal("failed to read json string from bytes", err)
	}
	if arb["key"] != expected {
		t.Fatalf("unexpected value for key %s, expected %s", arb["key"], expected)
	}
}

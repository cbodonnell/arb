package arb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Arb map[string]interface{}

// Create a new Arb
func New() Arb {
	return make(Arb)
}

// Read Arb
func Read(r io.Reader) (Arb, error) {
	decoder := json.NewDecoder(r)
	var a Arb
	err := decoder.Decode(&a)
	return a, err
}

// Write Arb
func (a Arb) Write(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err := e.Encode(a)
	return err
}

// Convert Arb to bytes
func (a Arb) ToBytes() []byte {
	var buf bytes.Buffer
	a.Write(&buf)
	return buf.Bytes()
}

// Convert Arb to string
func (a Arb) ToString() string {
	var buf bytes.Buffer
	a.Write(&buf)
	return buf.String()
}

// Get type of a property
func (a Arb) GetType(prop string) interface{} {
	p := a[prop]
	switch t := p.(type) {
	default:
		return t
	}
}

// Check if property exists
func (a Arb) Exists(prop string) bool {
	return a[prop] != nil
}

// Check if property is bool
func (a Arb) IsBool(prop string) bool {
	_, r := a[prop].(bool)
	return r
}

// Check if property is number
func (a Arb) IsNumber(prop string) bool {
	_, r := a[prop].(float64)
	return r
}

// Check if property is string
func (a Arb) IsString(prop string) bool {
	_, r := a[prop].(string)
	return r
}

// Check if property is an array
func (a Arb) IsArray(prop string) bool {
	_, r := a[prop].([]interface{})
	return r
}

// Check if property is an Arb
func (a Arb) IsArb(prop string) bool {
	_, r := a[prop].(map[string]interface{})
	return r
}

// Check if property is a URL
func (a Arb) IsURL(prop string) bool {
	s, err := a.GetString(prop)
	if err != nil {
		return false
	}
	_, err = url.Parse(s)
	if err != nil {
		return false
	}
	return true
}
func (a Arb) GetBool(prop string) (bool, error) {
	if s, ok := a[prop].(bool); !ok {
		return s, errors.New(fmt.Sprintf("%s is not a bool", prop))
	} else {
		return s, nil
	}
}

func (a Arb) GetNumber(prop string) (float64, error) {
	if s, ok := a[prop].(float64); !ok {
		return s, errors.New(fmt.Sprintf("%s is not a number", prop))
	} else {
		return s, nil
	}
}

func (a Arb) GetString(prop string) (string, error) {
	if s, ok := a[prop].(string); !ok {
		return s, errors.New(fmt.Sprintf("%s is not a string", prop))
	} else {
		return s, nil
	}
}

func (a Arb) GetArray(prop string) ([]interface{}, error) {
	if s, ok := a[prop].([]interface{}); !ok {
		return s, errors.New(fmt.Sprintf("%s is not an array", prop))
	} else {
		return s, nil
	}
}

func (a Arb) GetArb(prop string) (Arb, error) {
	if m, ok := a[prop].(map[string]interface{}); !ok {
		if s, ok := a[prop].(Arb); !ok {
			return s, errors.New(fmt.Sprintf("%s is not an Arb", prop))
		} else {
			return s, nil
		}
	} else {
		return m, nil
	}
}

func (a Arb) GetURL(prop string) (*url.URL, error) {
	s, err := a.GetString(prop)
	if err != nil {
		return nil, err
	}
	iri, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return iri, nil
}

// Convert a string property to an array
func (a Arb) PropToArray(prop string) error {
	if !a.IsArray(prop) {
		s, err := a.GetString(prop)
		if err != nil {
			return err
		}
		a[prop] = []string{s}
	}
	return nil
}

// Get Arb even if it is an IRI
// TODO: Allow headers
func (a Arb) FindArb(prop string) (Arb, error) {
	iri, err := a.GetURL(prop)
	if err != nil {
		return a.GetArb(prop)
	}
	client := http.DefaultClient
	req, err := http.NewRequest("GET", iri.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	arb, err := Read(resp.Body)
	if err != nil {
		return nil, err
	}
	return arb, nil
}

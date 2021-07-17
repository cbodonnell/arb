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

func New() Arb {
	return make(Arb)
}

func Read(r io.Reader) (Arb, error) {
	decoder := json.NewDecoder(r)
	var a Arb
	err := decoder.Decode(&a)
	return a, err
}

func (a Arb) Write(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err := e.Encode(a)
	return err
}

func (a Arb) ToBytes() []byte {
	var buf bytes.Buffer
	a.Write(&buf)
	return buf.Bytes()
}

func (a Arb) ToString() string {
	var buf bytes.Buffer
	a.Write(&buf)
	return buf.String()
}

func (a Arb) GetType(prop string) interface{} {
	p := a[prop]
	switch t := p.(type) {
	default:
		return t
	}
}

func (a Arb) Exists(prop string) bool {
	return a[prop] != nil
}

func (a Arb) IsBool(prop string) bool {
	_, r := a[prop].(bool)
	return r
}

func (a Arb) IsNumber(prop string) bool {
	_, r := a[prop].(float64)
	return r
}

func (a Arb) IsString(prop string) bool {
	_, r := a[prop].(string)
	return r
}

func (a Arb) IsArray(prop string) bool {
	_, r := a[prop].([]interface{})
	return r
}

func (a Arb) IsArb(prop string) bool {
	_, r := a[prop].(Arb)
	return r
}

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
	if s, ok := a[prop].(Arb); !ok {
		return s, errors.New(fmt.Sprintf("%s is not an Arb", prop))
	} else {
		return s, nil
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

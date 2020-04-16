package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonDocument struct {
	blob interface{}
}

func NewJsonDocument(content string) (*JsonDocument, error) {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, fmt.Errorf("Expected json string, got an empty string\n")
	}
	if content[0] == '{' { // out is object
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(content), &obj)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: obj}, nil
	} else if content[0] == '[' { // out is array
		arr := make([]interface{}, 0)
		err := json.Unmarshal([]byte(content), &arr)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: arr}, nil
	} else {
		return nil, fmt.Errorf("Expected [ or { as the json's first token, got \"%c\"\n", content[0])
	}
}

type JsonQuery struct {
	doc *JsonDocument
}

func NewQuery(doc *JsonDocument) *JsonQuery {
	return &JsonQuery{doc: doc}
}

func (j *JsonQuery) Select(tokens ...interface{}) (interface{}, error) {
	return rquery(j.doc.blob, tokens...)
}

func rquery(blob interface{}, tokens ...interface{}) (interface{}, error) {
	val := blob
	var err error
	for _, token := range tokens {
		val, err = query(val, token)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

func query(blob interface{}, query interface{}) (interface{}, error) {
	idx, ok := query.(int)
	if ok { // array
		arr, ok := blob.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Array index on non-array %v\n", blob)
		}
		if len(arr) <= idx { // out of bounds
			return nil, fmt.Errorf("Array index %d on array %v out of bounds\n", idx, blob)
		}
		return arr[idx], nil
	}
	tok, ok := query.(string)
	if ok { // object
		obj, ok := blob.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Object lookup \"%s\" on non-object %v\n", query, blob)
		}
		val, ok := obj[tok]
		if !ok { // field not exist
			return nil, fmt.Errorf("Object %v does not contain field %s\n", blob, query)
		}
		return val, nil
	}
	return nil, fmt.Errorf("Input %v is a non-array ans non-object\n", blob)
}

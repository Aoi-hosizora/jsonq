package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonDocument struct {
	blob interface{} // map / array
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

type MultiToken struct {
	sels []interface{}
}

func NewMultiToken(tokens ...interface{}) *MultiToken {
	return &MultiToken{sels: tokens}
}

func (j *JsonQuery) Select(tokens ...interface{}) (interface{}, error) {
	vals, multi, err := rquery(j.doc.blob, tokens...)
	if err != nil {
		return nil, err
	}
	if multi {
		return vals, nil
	}
	return vals[0], nil
}

func rquery(blob interface{}, tokens ...interface{}) ([]interface{}, bool, error) {
	vals := []interface{}{blob}
	isArray := false
	for _, token := range tokens { // take token of layers
		mtok, ok := token.(*MultiToken)
		if !ok {
			// current layer is single token
			for idx, val := range vals { // for all data
				val, err := query(val, token)
				if err != nil {
					return nil, isArray, err
				}
				vals[idx] = val // replace value directly
			}
		} else {
			// current layer is multi token
			isArray = true
			tmpVal := make([]interface{}, 0)

			for _, val := range vals { // for all data
				for _, stok := range mtok.sels {
					// for a token in multiToken (same level first)
					val, err := query(val, stok)
					if err != nil {
						return nil, isArray, err
					}
					tmpVal = append(tmpVal, val) // append to a new value array
				}
			}
			vals = tmpVal
		}
	}
	return vals, isArray, nil
}

func query(blob interface{}, token interface{}) (interface{}, error) {
	idx, ok := token.(int)
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
	tok, ok := token.(string)
	if ok { // object
		obj, ok := blob.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Object lookup \"%s\" on non-object %v\n", token, blob)
		}
		val, ok := obj[tok]
		if !ok { // field not exist
			return nil, fmt.Errorf("Object %v does not contain field \"%s\"\n", blob, token)
		}
		return val, nil
	}
	return nil, fmt.Errorf("Input %v is a non-array and non-object\n", blob)
}

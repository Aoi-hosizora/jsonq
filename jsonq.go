package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
)

// blob is a interface of map[string]interface{} (if it is an object-wrapped json)
// or []interface{} (if it is an array-wrapped json)
type JsonDocument struct {
	blob interface{}
}

// create a JsonDocument, parse json string first
func NewJsonDocument(content string) (*JsonDocument, error) {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, fmt.Errorf("Expected json string, got an empty string\n")
	}
	if content[0] == '{' { // object-wrapped
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(content), &obj)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: obj}, nil
	} else if content[0] == '[' { // array-wrapped
		arr := make([]interface{}, 0)
		err := json.Unmarshal([]byte(content), &arr)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: arr}, nil
	} else { // other start token
		return nil, fmt.Errorf("Expected [ or { as the json's first token, got \"%c\"\n", content[0])
	}
}

// doc is a json document that has been parse correctly
type JsonQuery struct {
	doc *JsonDocument
}

// create a json query to select json
func NewQuery(doc *JsonDocument) *JsonQuery {
	return &JsonQuery{doc: doc}
}

// select multiple fields in the same layer
type MultiToken struct {
	sels []interface{}
}

// build a selector which will select multiple fields in the same layer
func NewMultiToken(tokens ...interface{}) *MultiToken {
	return &MultiToken{sels: tokens}
}

// key code start from here

// select json by a slice of strings / integers / MultiTokens
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

// select json by a selector string
func (j *JsonQuery) SelectBySelector(selectorString string) (interface{}, error) {
	selector, err := escapeSelector(selectorString)
	if err != nil {
		return nil, err
	}
	return j.Select(selector...)
}

// repetition query: tokens []interface{}
// If it is a singleToken(string, integer), it will select fields in different layers
// If it is a multiToken, it will select fields in the same layer
// once have a MultiToken, that will return an array
func rquery(blob interface{}, tokens ...interface{}) ([]interface{}, bool, error) {
	vals := []interface{}{blob}
	isArray := false
	for _, token := range tokens {
		// get a token (stok / mtok) in different layers
		mtok, ok := token.(*MultiToken)

		if !ok {
			// current layer is single token
			for idx, val := range vals { // for all data
				val, err := query(val, token)
				if err != nil {
					return nil, isArray, err
				}
				vals[idx] = val // replace values directly
			}
		} else {
			// current layer is multi token
			isArray = true
			tmpVal := make([]interface{}, 0)
			for _, val := range vals { // for all data in the current array
				for _, stok := range mtok.sels {
					// get a single token in mtok (same layer first)
					val, err := query(val, stok)
					if err != nil {
						return nil, isArray, err
					}
					tmpVal = append(tmpVal, val) // append to a new value array
				}
			}
			vals = tmpVal // replace values entirely
		}
	}
	return vals, isArray, nil
}

// single query: token interface{}
// If it is an integer, it will select an item in the array
// If it is a string, it will select a field in the map
// If the index is out of bound, or the map does not contain field, it will return an error
func query(blob interface{}, token interface{}) (interface{}, error) {
	idx, ok := token.(int) // index
	if ok {
		arr, ok := blob.([]interface{}) // array
		if !ok {
			return nil, fmt.Errorf("Array index on non-array %v\n", blob)
		}
		if len(arr) <= idx { // out of bound
			return nil, fmt.Errorf("Array index %d on array %v out of bound\n", idx, blob)
		}
		return arr[idx], nil
	}

	tok, ok := token.(string) // key
	if ok {
		obj, ok := blob.(map[string]interface{}) // object
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

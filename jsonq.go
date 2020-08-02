package jsonq

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Parse json string first for json query.
type JsonDocument struct {
	// a interface of
	// 1. `map[string]interface{}` (if it is an object-wrapped json)
	// 2. `[]interface{}` (if it is an array-wrapped json)
	blob interface{}
}

// Create a JsonDocument, handle json string first (`{` or `[`).
func NewJsonDocument(data []byte) (*JsonDocument, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, fmt.Errorf("expected json string, got an empty string")
	}

	// object-wrapped
	if data[0] == '{' {
		obj := make(map[string]interface{})
		err := json.Unmarshal(data, &obj)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: obj}, nil
	}

	// array-wrapped
	if data[0] == '[' {
		arr := make([]interface{}, 0)
		err := json.Unmarshal(data, &arr)
		if err != nil {
			return nil, err
		}
		return &JsonDocument{blob: arr}, nil
	}

	// other start token
	return nil, fmt.Errorf("expected [ or { as the json's first token, got \"%c\"", data[0])
}

// Query json fields.
type JsonQuery struct {
	// a json document that has been check (parse) correctly
	doc *JsonDocument
}

// Create a JsonQuery to query json.
func NewJsonQuery(doc *JsonDocument) *JsonQuery {
	return &JsonQuery{doc: doc}
}

// Select multiple fields in the same layer -> "+".
type multiToken struct {
	sels []interface{}
}

// Build a multiple selector which will select multiple fields in the same layer.
func Multi(tokens ...interface{}) *multiToken {
	return &multiToken{sels: tokens}
}

// Select all fields in the same layer -> "*".
type starToken struct{}

// Build a selector which will select all fields in the same layer.
func All() *starToken {
	return &starToken{}
}

// ========================
// key code start from here
// ========================

// Query json by a slice of strings / integers / MultiTokens.
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

// Query json by a selector string.
func (j *JsonQuery) SelectBySelector(selectorString string) (interface{}, error) {
	selector, err := _NewParser(selectorString).Parse()
	if err != nil {
		return nil, err
	}
	return j.Select(selector...)
}

// Repetition query: tokens []interface{}.
//
// If it is a SingleToken(string, integer), it will select fields in different layers.
// If it is a multiToken(starToken will return error), it will select fields in the same layer.
// If it is a starToken, it will select all fields in the same layer.
// Once have a multiToken or an starToken, that will return an array.
func rquery(blob interface{}, tokens ...interface{}) ([]interface{}, bool, error) {
	vals := []interface{}{blob}
	isArray := false
	for _, token := range tokens {
		// get a token (stok / mtok / atok) in different layers
		mtok, isMul := token.(*multiToken)
		_, isAll := token.(*starToken)

		if !isMul && !isAll {
			// current layer is a single token
			for idx, val := range vals { // for all data
				val, err := query(val, token)
				if err != nil {
					return nil, isArray, err
				}
				vals[idx] = val // replace values directly
			}
		} else {
			// current layer is a multi token / an star token
			isArray = true
			tmpVal := make([]interface{}, 0)

			// for all data in the current array
			if isMul {
				// current layer is a multi token
				for _, val := range vals {
					for _, stok := range mtok.sels {
						// get a single token in mtok (same layer first)
						val, err := query(val, stok)
						if err != nil {
							return nil, isArray, err
						}
						tmpVal = append(tmpVal, val) // append to a new value array
					}
				}
			} else {
				// current layer is a star token
				for _, val := range vals {
					// get all fields
					vals, err := queryAll(val)
					if err != nil {
						return nil, isArray, err
					}
					tmpVal = append(tmpVal, vals...) // append to a new value array
				}
			}

			vals = tmpVal // replace values entirely
		}
	}
	return vals, isArray, nil
}

// Query a single field: token interface{}.
//
// If it is an integer, it will select an item in the array.
// If it is a string, it will select a field in the map.
// If the index is out of bound, or the map does not contain field, it will return an error.
func query(blob interface{}, token interface{}) (interface{}, error) {
	idx, ok := token.(int) // index
	if ok {
		arr, ok := blob.([]interface{}) // array
		if !ok {
			return nil, fmt.Errorf("Array index on non-array %v\n", blob)
		}
		if len(arr) <= idx || idx <= -len(arr)-1 { // out of bound
			return nil, fmt.Errorf("Array index %d on array %v out of bound\n", idx, blob)
		}
		if idx < 0 {
			idx += len(arr)
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

// Query all fields: starToken.
func queryAll(blob interface{}) ([]interface{}, error) {
	arr, ok := blob.([]interface{})
	if ok {
		return arr, nil
	}

	obj, ok := blob.(map[string]interface{})
	if ok {
		out := make([]interface{}, len(obj))
		idx := 0
		for k := range obj {
			out[idx] = obj[k]
			idx++
		}
		return out, nil
	}

	return nil, fmt.Errorf("Input %v is a non-array and non-object\n", blob)
}

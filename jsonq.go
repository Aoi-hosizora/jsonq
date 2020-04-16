package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonDocument struct {
	obj map[string]interface{}
}

func NewJsonDocument(content string) (*JsonDocument, error) {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, fmt.Errorf("")
	}
	obj := make(map[string]interface{})
	if content[0] == '{' {
		err := json.Unmarshal([]byte(content), &obj)
		if err != nil {
			return nil, err
		}
	} else if content[0] == '[' {
		arr := make([]interface{}, 0)
		err := json.Unmarshal([]byte(content), &arr)
		if err != nil {
			return nil, err
		}
		obj["_"] = arr
	} else {
		return nil, fmt.Errorf("")
	}
	return &JsonDocument{obj: obj}, nil
}

type JsonQuery struct {
	doc *JsonDocument
}

func NewQuery(doc *JsonDocument) *JsonQuery {
	return &JsonQuery{doc: doc}
}

func (j *JsonQuery) Select(tokens ...interface{}) (interface{}, error) {
	if len(tokens) == 0 {
		return j.doc.obj, nil
	}
	idx, ok := tokens[0].(int)
	if ok { // the first one is an integer -> out is array
		arr, ok := j.doc.obj["_"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("")
		}
		if len(arr) <= idx {
			return nil, fmt.Errorf("")
		}
		return rquery(arr[idx], tokens[1:]...)
	}
	return rquery(j.doc.obj, tokens...)
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
			return nil, fmt.Errorf("")
		}
		if len(arr) <= idx { // out of bound
			return nil, fmt.Errorf("")
		}
		return arr[idx], nil
	}
	tok, ok := query.(string)
	if ok { // object
		obj, ok := blob.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("")
		}
		val, ok := obj[tok]
		if !ok { // key not exist
			return nil, fmt.Errorf("")
		}
		return val, nil
	}
	return nil, fmt.Errorf("")
}

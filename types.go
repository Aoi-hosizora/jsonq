package jsonq

import (
	"fmt"
)

func interfaceToBool(i interface{}) (bool, error) {
	if b, ok := i.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("Excepted a bool value, got \"%v\"\n", i)
}

func interfaceToInt64(i interface{}) (int64, error) {
	switch i.(type) {
	case int64:
		return i.(int64), nil
	case float64:
		return int64(i.(float64)), nil
	}
	return 0, fmt.Errorf("Excepted an int64 value, got \"%v\"\n", i)
}

func interfaceToFloat64(i interface{}) (float64, error) {
	switch i.(type) {
	case float64:
		return i.(float64), nil
	case int64:
		return float64(i.(int64)), nil
	}
	return 0, fmt.Errorf("Excepted a float64 value, got \"%v\"\n", i)
}

func interfaceToString(i interface{}) (string, error) {
	if b, ok := i.(string); ok {
		return b, nil
	}
	return "", fmt.Errorf("Excepted a string value, got \"%v\"\n", i)
}

func interfaceToObject(i interface{}) (map[string]interface{}, error) {
	if b, ok := i.(map[string]interface{}); ok {
		return b, nil
	}
	return nil, fmt.Errorf("Excepted an object value, got \"%v\"\n", i)
}

func interfaceToArray(i interface{}) ([]interface{}, error) {
	if b, ok := i.([]interface{}); ok {
		return b, nil
	}
	return nil, fmt.Errorf("Excepted an array value, got \"%v\"\n", i)
}

// ===========================================================================

func (j *JsonQuery) Bool(tokens ...interface{}) (bool, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return false, err
	}
	return interfaceToBool(res)
}

func (j *JsonQuery) BoolBySelector(selectorString string) (bool, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return false, err
	}
	return interfaceToBool(res)
}

func (j *JsonQuery) Int64(tokens ...interface{}) (int64, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return 0, err
	}
	return interfaceToInt64(res)
}

func (j *JsonQuery) Int64BySelector(selectorString string) (int64, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return 0, err
	}
	return interfaceToInt64(res)
}

func (j *JsonQuery) Float64(tokens ...interface{}) (float64, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return 0, err
	}
	return interfaceToFloat64(res)
}

func (j *JsonQuery) Float64BySelector(selectorString string) (float64, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return 0, err
	}
	return interfaceToFloat64(res)
}

func (j *JsonQuery) String(tokens ...interface{}) (string, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return "", err
	}
	return interfaceToString(res)
}

func (j *JsonQuery) StringBySelector(selectorString string) (string, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return "", err
	}
	return interfaceToString(res)
}

func (j *JsonQuery) Object(tokens ...interface{}) (map[string]interface{}, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	return interfaceToObject(res)
}

func (j *JsonQuery) ObjectBySelector(selectorString string) (map[string]interface{}, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return nil, err
	}
	return interfaceToObject(res)
}

func (j *JsonQuery) Array(tokens ...interface{}) ([]interface{}, error) {
	res, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	return interfaceToArray(res)
}

func (j *JsonQuery) ArrayBySelector(selectorString string) ([]interface{}, error) {
	res, err := j.SelectBySelector(selectorString)
	if err != nil {
		return nil, err
	}
	return interfaceToArray(res)
}

// ===========================================================================

func (j *JsonQuery) Bools(tokens ...interface{}) ([]bool, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]bool, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToBool(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) BoolsBySelector(selector string) ([]bool, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]bool, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToBool(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Int64s(tokens ...interface{}) ([]int64, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToInt64(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Int64sBySelector(selector string) ([]int64, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToInt64(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Float64s(tokens ...interface{}) ([]float64, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]float64, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToFloat64(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Float64sBySelector(selector string) ([]float64, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]float64, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToFloat64(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Strings(tokens ...interface{}) ([]string, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToString(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) StringsBySelector(selector string) ([]string, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToString(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Objects(tokens ...interface{}) ([]map[string]interface{}, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]map[string]interface{}, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToObject(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) ObjectsBySelector(selector string) ([]map[string]interface{}, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([]map[string]interface{}, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToObject(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) Arrays(tokens ...interface{}) ([][]interface{}, error) {
	itf, err := j.Select(tokens...)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([][]interface{}, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToArray(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *JsonQuery) ArraysBySelector(selector string) ([][]interface{}, error) {
	itf, err := j.SelectBySelector(selector)
	if err != nil {
		return nil, err
	}
	arr, err := interfaceToArray(itf)
	if err != nil {
		return nil, err
	}
	res := make([][]interface{}, len(arr))
	for idx := range arr {
		res[idx], err = interfaceToArray(arr[idx])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

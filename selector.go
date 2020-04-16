package jsonq

import (
	"fmt"
	"strconv"
	"strings"
)

func renderSelector(selector string) ([]interface{}, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return []interface{}{}, nil
	}

	// selector string -> [][]string
	_allToks := make([][]string, 0)
	_curToks := make([]string, 0)
	_curTok := ""

	for idx := 0; idx <= len(selector); idx++ {
		if idx == len(selector) { // finish
			_allToks = append(_allToks, append(_curToks, _curTok)) // last MultiToken
			break
		}
		if selector[idx] == ' ' { // new token
			_allToks = append(_allToks, append(_curToks, _curTok)) // last tokenString
			_curToks = make([]string, 0)
			_curTok = ""
			continue
		}
		if selector[idx] == '+' { // new multiToken
			_curToks = append(_curToks, _curTok)
			_curTok = ""
			continue
		}

		ap := string(selector[idx])
		if ap == "\\" && len(selector) > idx { // escape all char, including " ", "+" and "\"
			idx++
			ap = string(selector[idx])
		}
		_curTok += ap // append to current token string
	}

	// remove empty token
	allToks := make([][]string, 0)
	curToks := make([]string, 0)
	for _, mtok := range _allToks {
		for _, stok := range mtok {
			if stok != "" {
				curToks = append(curToks, stok)
			}
		}
		if len(curToks) != 0 {
			allToks = append(allToks, curToks)
			curToks = make([]string, 0)
		}
	}

	// selector [][]string -> []interface{}

	ret := make([]interface{}, 0)
	for _, mtok := range allToks {
		if len(mtok) == 1 { // single token
			err := render(&ret, mtok[0])
			if err != nil {
				return nil, err
			}
		} else if len(mtok) >= 2 { // multi token
			multi := NewMultiToken()
			for _, stok := range mtok { // for each single token
				err := render(&multi.sels, stok)
				if err != nil {
					return nil, err
				}
			}
			ret = append(ret, multi)
		}
	}

	return ret, nil
}

// render number and string
func render(arr *[]interface{}, tok string) error {
	if len(tok) == 0 {
		return nil
	}

	// len(tok) >= 1
	if tok[0] != '#' { // xxx
		*arr = append(*arr, tok)
		return nil
	}

	// tok[0] == '#' && len(tok) >= 1
	if len(tok) == 1 { // #
		return fmt.Errorf("Number should appear after single #, find null\n")
	}

	// tok[0] == '#' && len(tok) >= 2
	if tok[1] == '#' { // ##xxx
		*arr = append(*arr, tok[1:]) // delete one #
		return nil
	}

	// tok[0] == '#' && tok[1] != '#'
	number, err := strconv.Atoi(tok[1:])
	if err != nil { // #xxx
		return fmt.Errorf("Number should appear after single #, find \"%c\"\n", tok[1])
	}

	// #0
	*arr = append(*arr, number)
	return nil
}

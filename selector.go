package jsonq

import (
	"fmt"
	"strconv"
	"strings"
)

// render a selector string to array of interface{} (string || int || MultiToken)
func escapeSelector(selector string) ([]interface{}, error) {
	allToks := escapeString(selector) // [][]string
	ret := make([]interface{}, 0)
	for _, mtok := range allToks {
		if len(mtok) == 1 { // single token / star token
			res, err := escapeSharp(mtok[0]) // string || int || star
			if err != nil {
				return nil, err
			}
			if res != nil {
				ret = append(ret, res) // <<<
			}
		} else if len(mtok) >= 2 { // multi token
			multi := NewMultiToken()
			for _, stok := range mtok { // for each single token
				res, err := escapeSharp(stok) // string || int
				if err != nil {
					return nil, err
				}
				if res != nil {
					multi.sels = append(multi.sels, res) // <<<
				}
			}
			if len(multi.sels) != 0 {
				ret = append(ret, multi) // <<<
			}
		}
	}

	return ret, nil
}

// render "\" "+" " " in a selector string, return a multiTok array
func escapeString(selector string) [][]string {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return nil
	}

	// string to [][]string
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
		if ap == "\\" && len(selector) > idx { // escape all char, including " ", "+" "\" "#" "*"
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

	return allToks
}

// render # in token, return number or string or star
func escapeSharp(tok string) (interface{}, error) {
	if len(tok) == 0 {
		return nil, nil // return nil
	}
	if tok == "*" { // *
		return NewStarToken(), nil
	}
	if len(tok) >= 2 && tok[0] == '*' { // *xx
		return tok[1:], nil
	}
	if tok[0] != '#' { // xxx
		return tok, nil
	}

	// start #
	// tok[0] == '#' && len(tok) >= 1
	if len(tok) == 1 { // #
		return nil, fmt.Errorf("Number should appear after single #, find null\n")
	}

	// tok[0] == '#' && len(tok) >= 2
	if tok[1] == '#' { // ##xxx
		return tok[1:], nil // delete one #
	}

	// tok[0] == '#' && tok[1] != '#'
	number, err := strconv.Atoi(tok[1:])
	if err != nil { // #xxx
		return nil, fmt.Errorf("Number should appear after single #, find \"%c\"\n", tok[1])
	}

	// #0
	return number, nil
}

package jsonq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type _Token int

const (
	eof = rune(0)

	_ILLEGAL _Token = iota
	_EOF
	_WHITESPACE
	_IDENT
	_ASTERISK // *
	_PLUS     // +
	_NUMBER   // #0
)

type _Scanner struct {
	r *bufio.Reader
}

func _NewScanner(r io.Reader) *_Scanner {
	return &_Scanner{r: bufio.NewReader(r)}
}

func (s *_Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if ch == 0 || err != nil {
		return eof
	}
	return ch
}

func (s *_Scanner) unread() {
	_ = s.r.UnreadByte()
}

func (s *_Scanner) Scan() (tok _Token, lit string, err error) {
	ch := s.read()
	if isWhitespace(ch) { // -> next layer
		s.unread() // release the previous ws
		return s.scanWhitespace()
	} else if isSharp(ch) { // -> number (include -)
		return s.scanNumber() // start with #
	} else if isStar(ch) { // -> all fields
		return s.scanStar() // start with and only be *
	} else if isIdent(ch) { // -> string (include \)
		s.unread() // release the previous char
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return _EOF, "", nil
	case '+':
		return _PLUS, "+", nil // -> new fields
	default:
		return _ILLEGAL, "", fmt.Errorf("Illegal char as the start with selector\n")
	}
}

func (s *_Scanner) scanWhitespace() (tok _Token, lit string, err error) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return _WHITESPACE, " ", nil
}

func (s *_Scanner) scanNumber() (tok _Token, lit string, err error) {
	var buf bytes.Buffer
	minus := false
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isWhitespace(ch) || isPlus(ch) { // next layer or next field
			s.unread()
			break
		} else if isMinus(ch) {
			if minus || len(buf.String()) != 0 {
				return _ILLEGAL, "", fmt.Errorf("Could mix number and string after #\n")
			} else {
				minus = true
				buf.WriteRune(ch)
			}
		} else if isDigit(ch) {
			buf.WriteRune(ch)
		} else {
			return _ILLEGAL, "", fmt.Errorf("Could mix number and string after #\n")
		}
	}

	if buf.String() == "" {
		return _NUMBER, "0", nil
	} else {
		return _NUMBER, buf.String(), nil
	}
}

func (s *_Scanner) scanStar() (tok _Token, lit string, err error) {
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isWhitespace(ch) { // next layer
			s.unread()
			break
		} else if isPlus(ch) { // next field
			return _ILLEGAL, "", fmt.Errorf("Could not select the next field when use *\n")
		} else {
			return _ILLEGAL, "", fmt.Errorf("Could not mix * and other token after *\n")
		}
	}
	return _ASTERISK, "*", nil
}

func (s *_Scanner) scanIdent() (tok _Token, lit string, err error) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isWhitespace(ch) || isPlus(ch) { // next layer or next field
			s.unread()
			break
		} else if isBackSlash(ch) { // escape (specially when start with # * and contain ws +)
			ch2 := s.read()
			if ch2 == eof {
				break
			}
			buf.WriteRune(ch2)
		} else {
			buf.WriteRune(ch)
		}
	}
	return _IDENT, buf.String(), nil
}

type _Parser struct {
	s *_Scanner
}

func _NewParser(r string) *_Parser {
	return &_Parser{s: _NewScanner(strings.NewReader(r))}
}

func (p *_Parser) readNextTok() (tok _Token, lit string, err error) {
	return p.s.Scan()
}

func (p *_Parser) Parse() (selector []interface{}, err error) {
	toks := []*multiToken{{}} // number / string / multiToken / starToken

out:
	for {
		tok, lit, err := p.readNextTok()
		if err != nil {
			return nil, err
		}

		switch tok {
		case _EOF:
			break out
		case _WHITESPACE:
			toks = append(toks, &multiToken{})
		case _PLUS: // -> no need to handle, append to the last mtok directly
		case _NUMBER:
			num, err := strconv.Atoi(lit)
			if err != nil {
				panic(err)
			}
			toks[len(toks)-1].sels = append(toks[len(toks)-1].sels, num)
		case _ASTERISK:
			toks[len(toks)-1].sels = append(toks[len(toks)-1].sels, All())
		case _IDENT:
			toks[len(toks)-1].sels = append(toks[len(toks)-1].sels, lit)
		default:
			panic("Illegal token type\n")
		}
	}

	out := make([]interface{}, 0)
	for _, mtok := range toks {
		if len(mtok.sels) == 0 {
			continue
		} else if len(mtok.sels) == 1 {
			out = append(out, mtok.sels[0])
		} else {
			out = append(out, mtok)
		}
	}
	return out, nil
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isPlus(ch rune) bool {
	return ch == '+'
}

func isSharp(ch rune) bool {
	return ch == '#'
}

func isStar(ch rune) bool {
	return ch == '*'
}

func isIdent(ch rune) bool {
	return ch != '#' && ch != '*' && ch != ' ' && ch != '+' && ch != eof
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isMinus(ch rune) bool {
	return ch == '-'
}

func isBackSlash(ch rune) bool {
	return ch == '\\'
}

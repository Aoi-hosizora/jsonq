package jsonq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Token int

const (
	eof = rune(0)

	ILLEGAL Token = iota
	EOF
	WHITESPACE
	IDENT
	ASTERISK // *
	PLUS     // +
	NUMBER   // #0
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if ch == 0 || err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadByte()
}

func (s *Scanner) Scan() (tok Token, lit string, err error) {
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
		return EOF, "", nil
	case '+':
		return PLUS, "+", nil // -> new fields
	default:
		return ILLEGAL, "", fmt.Errorf("Illegal char as the start with selector\n")
	}
}

func (s *Scanner) scanWhitespace() (tok Token, lit string, err error) {
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
	return WHITESPACE, " ", nil
}

func (s *Scanner) scanNumber() (tok Token, lit string, err error) {
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
				return ILLEGAL, "", fmt.Errorf("Could mix number and string after #\n")
			} else {
				minus = true
				buf.WriteRune(ch)
			}
		} else if isDigit(ch) {
			buf.WriteRune(ch)
		} else {
			return ILLEGAL, "", fmt.Errorf("Could mix number and string after #\n")
		}
	}

	if buf.String() == "" {
		return NUMBER, "0", nil
	} else {
		return NUMBER, buf.String(), nil
	}
}

func (s *Scanner) scanStar() (tok Token, lit string, err error) {
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isWhitespace(ch) { // next layer
			s.unread()
			break
		} else if isPlus(ch) { // next field
			return ILLEGAL, "", fmt.Errorf("Could not select the next field when use *\n")
		} else {
			return ILLEGAL, "", fmt.Errorf("Could not mix * and other token after *\n")
		}
	}
	return ASTERISK, "*", nil
}

func (s *Scanner) scanIdent() (tok Token, lit string, err error) {
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
	return IDENT, buf.String(), nil
}

type Parser struct {
	s *Scanner
}

func NewParser(r string) *Parser {
	return &Parser{s: NewScanner(strings.NewReader(r))}
}

func (p *Parser) readNextTok() (tok Token, lit string, err error) {
	return p.s.Scan()
}

func (p *Parser) Parse() (selector []interface{}, err error) {
	toks := []*MultiToken{{}} // number / string / multiToken / starToken

out:
	for {
		tok, lit, err := p.readNextTok()
		if err != nil {
			return nil, err
		}

		switch tok {
		case EOF:
			break out
		case WHITESPACE:
			toks = append(toks, &MultiToken{})
		case PLUS: // -> no need to handle, append to the last mtok directly
		case NUMBER:
			num, err := strconv.Atoi(lit)
			if err != nil {
				panic(err)
			}
			toks[len(toks)-1].sels = append(toks[len(toks)-1].sels, num)
		case ASTERISK:
			toks[len(toks)-1].sels = append(toks[len(toks)-1].sels, NewStarToken())
		case IDENT:
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

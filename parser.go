package jsonq

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	eof = rune(0)

	EOF Token = iota
	WHITESPACE
	IDENT    // (#|*)?xxx
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
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadByte()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isBackSlash(ch rune) bool {
	return ch == '\\'
}

func isSharp(ch rune) bool {
	return ch == '#'
}

func isAsterisk(ch rune) bool {
	return ch == '*'
}

func isPlus(ch rune) bool {
	return ch == '+'
}

func isDigit(ch rune) bool {
	return ch == '-' || (ch >= '0' && ch <= '9')
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isSharp(ch) {
		s.unread()
		return s.scanIdent(true)
	} else if isAsterisk(ch) {
		return s.scanIdent(false)
	}

	switch ch {
	case eof:
		return EOF, ""
	default:
		return s.scanIdent(false)
	}
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return WHITESPACE, buf.String()
}

func (s *Scanner) scanIdent(startWithSharp bool) (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	tok = EOF
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isBackSlash(ch) {
			ch2 := s.read()
			_, _ = buf.WriteRune(ch2)
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else if startWithSharp && isDigit(ch) {
			if tok == EOF {
				tok = NUMBER
			} else if tok != NUMBER {
				panic("Could not mix number and string when starts with #\n")
			}
			_, _ = buf.WriteRune(ch)
		} else {
			if tok == NUMBER {
				panic("Could not mix number and string when starts with #\n")
			}
			_, _ = buf.WriteRune(ch)
		}
	}

	if buf.String() == "*" {
		return ASTERISK, "*"
	} else if buf.String() == "#" {
		return NUMBER, "0"
	} else if tok == NUMBER {
		return NUMBER, buf.String()[1:]
	} else {
		return IDENT, buf.String()
	}
}

// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

const (
	_ = iota

	IDENT

	CHAR
	STRING

	INT_BIN
	INT_OCT
	INT_DEC
	INT_HEX

	LPAREN
	RPAREN
)

type Position struct {
	Offset, Line, Column int
}

func (p Position) String() string {
	return fmt.Sprint(p.Offset, ":", p.Line, ":", p.Column)
}

type BufferScanner struct {
	Position // Cursor
	Buffer   []rune
}

// Move returns current char and move cursor to the next.
// Move does not error when GetChar does not error.
func (bs *BufferScanner) Move() rune {
	ch := bs.GetChar()

	if ch == '\n' {
		bs.Column = 0
		bs.Line++
	} else {
		bs.Column++
	}

	bs.Offset++

	return ch
}

// GetChar returns the char at the cursor.
func (bs *BufferScanner) GetChar() rune {
	if bs.Offset == len(bs.Buffer) {
		panic(EOFError{Pos: bs.Position})
	}
	return bs.Buffer[bs.Offset]
}

// Scanner is the token scanner.
type Scanner struct {
	BufferScanner
}

func (s *Scanner) GotoNextLine() {
	for {
		ch := s.GetChar()
		if ch == '\n' {
			return
		}
		s.Move()
	}
}

func (s *Scanner) SkipWhitespace() {
	for {
		ch := s.GetChar()
		switch ch {
		case ' ':
		case '\t':
		case '\n':
		case '\r':
		default:
			return
		}
		s.Move()
	}
}

func (s *Scanner) ScanUnicodeCharHex(runesN int) rune {
	literal := make([]rune, runesN)
	for i := 0; i < runesN; i++ {
		ch := s.Move()
		literal[i] = ch
	}
	ch, err := strconv.ParseUint(string(literal), 16, runesN*4)
	if err != nil {
		switch {
		case errors.Is(err.(*strconv.NumError).Err, strconv.ErrRange):
			panic(FormatError{Pos: s.Position})
		case errors.Is(err.(*strconv.NumError).Err, strconv.ErrSyntax):
			panic(FormatError{Pos: s.Position})
		default:
			panic(err)
		}
	}
	return rune(ch)
}

// ScanEscapeChar returns the parsed char.
func (s *Scanner) ScanEscapeChar(quote rune) rune {
	ch := s.Move()
	switch ch {
	case quote:
		return quote
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case '\\':
		return '\\'
	case 'x': // Hex 1-byte unicode, 2 runes in total.
		return s.ScanUnicodeCharHex(2)
	case 'u': // Hex 2-byte unicode, 4 runes in total.
		return s.ScanUnicodeCharHex(4)
	case 'U': // Hex 4-byte unicode, 8 runes in total.
		return s.ScanUnicodeCharHex(8)
	default:
		panic(UnknownEscapeCharError{Char: ch})
	}
}

// ScanQuotedString scans the string or char.
// PANIC: Non-closed string might cause panic due to EOFError.
func (s *Scanner) ScanQuotedString(quote rune) (int, []rune) {
	s.Move()
	var str []rune
	for {
		ch := s.Move()
		switch ch {
		case '\\':
			ch := s.ScanEscapeChar(quote)
			str = append(str, ch)
		case quote:
			return STRING, str
		default:
			str = append(str, ch)
		}
	}
}

// ScanQuotedChar scans char.
func (s *Scanner) ScanQuotedChar() (int, []rune) {
	_, quote := s.ScanQuotedString('\'')
	if len(quote) != 1 {
		panic(FormatError{Pos: s.Position})
	}
	return CHAR, quote
}

func (s *Scanner) ScanWhile(cond func() bool) []rune {
	var seq []rune

	for cond() {
		seq = append(seq, s.GetChar())
		s.Move()
	}

	if len(seq) == 0 {
		panic(FormatError{Pos: s.Position})
	}

	return seq
}

func (s *Scanner) ScanBinDigit() (int, []rune) {
	cond := func() bool {
		ch := s.GetChar()
		return ch == '0' || ch == '1'
	}
	return INT_BIN, s.ScanWhile(cond)
}

func (s *Scanner) ScanOctDigit() (int, []rune) {
	cond := func() bool {
		ch := s.GetChar()
		return '0' <= ch && ch <= '7'
	}
	return INT_OCT, s.ScanWhile(cond)
}

func (s *Scanner) ScanHexDigit() (int, []rune) {
	cond := func() bool {
		ch := s.GetChar()
		return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f'
	}
	return INT_HEX, s.ScanWhile(cond)
}

func (s *Scanner) ScanDigit() (int, []rune) {
	ch := s.Move()
	digits := []rune{ch}

	if ch == '0' {
		switch s.Move() {
		case 'x':
			return s.ScanHexDigit()
		case 'o':
			return s.ScanOctDigit()
		case 'b':
			return s.ScanBinDigit()
		default:
			panic(FormatError{Pos: s.Position})
		}
	}

	for unicode.IsDigit(ch) {
		ch = s.Move()
		digits = append(digits, ch)
	}

	return INT_DEC, digits
}

// ScanWord scans and accepts only letters, digits and underlines.
// No valid string found when returns empty []rune.
func (s *Scanner) ScanWord() (int, []rune) {
	var word []rune
	for {
		ch := s.GetChar()
		switch {
		case unicode.IsDigit(ch):
		case unicode.IsLetter(ch):
		case ch == '_':
		default: // Terminate
			if len(word) == 0 {
				panic(FormatError{Pos: s.Position})
			}
			return IDENT, word
		}

		s.Move()
		word = append(word, ch)
	}
}

// ScanToken decides the next way to scan by the cursor.
func (s *Scanner) ScanToken() (Position, int, string) {
	s.SkipWhitespace()

	begin := s.Position

	kind, lit := func() (int, []rune) {
		ch := s.GetChar()

		switch {
		case unicode.IsDigit(ch): // Digital literal value
			return s.ScanDigit()

		case unicode.IsLetter(ch) || ch == '_': // Keyword OR Idents
			return s.ScanWord()

		default:
			switch ch {
			case '(':
				return LPAREN, []rune{s.Move()}

			case ')':
				return RPAREN, []rune{s.Move()}

			case '"': // String
				return s.ScanQuotedString(ch)

			case '\'': // Char
				return s.ScanQuotedChar()

			default:
				panic("impossible")
			}
		}
	}()

	return begin, kind, string(lit)
}

func ParseInt[T int | int8 | int16 | int32 | int64](kind int, lit string) T {
	parseInt := func(base int) T {
		i, err := strconv.ParseInt(lit, 16, 64)
		if err != nil {
			panic(err)
		}
		return T(i)
	}
	switch kind {
	case INT_HEX:
		return parseInt(16)
	case INT_DEC:
		return parseInt(10)
	case INT_OCT:
		return parseInt(8)
	case INT_BIN:
		return parseInt(2)
	default:
		panic("unexpected kind")
	}
}

func ParseUint[T uint | uint8 | uint16 | uint32 | uint64](kind int, lit string) T {
	parseUint := func(base int) T {
		i, err := strconv.ParseUint(lit, 16, 64)
		if err != nil {
			panic(err)
		}
		return T(i)
	}
	switch kind {
	case INT_HEX:
		return parseUint(16)
	case INT_DEC:
		return parseUint(10)
	case INT_OCT:
		return parseUint(8)
	case INT_BIN:
		return parseUint(2)
	default:
		panic("unexpected kind")
	}
}

//
// Copyright (C) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package lang

import (
	"fmt"
	"strconv"
)

type tokenType int

// Lexer token types
const (
	AOpenParen tokenType = iota
	ACloseParen
	ASymbol
	AString
	AInteger
	AFloat
	AQuote
)

// Token is the smallest unit of meaning for the little language
type Token struct {
	kind  tokenType
	value string
}

func (t Token) String() string {
	d := map[tokenType]string{
		AOpenParen:  "open-paren",
		ACloseParen: "close-paren",
		ASymbol:     "symbol",
		AString:     "string",
		AInteger:    "integer",
		AFloat:      "float",
		AQuote:      "quote",
	}
	return fmt.Sprintf("<#%v:[%+v]>", d[t.kind], t.value)
}

// Tokens contain the list of interpretable words in a form.
type Tokens struct {
	Tokens []Token
	word   []rune
	form   string
	kind   tokenType
}

func (ts *Tokens) pushChar(c rune) {
	ts.word = append(ts.word, c)
}

func (ts *Tokens) pushWord() {
	isFloat := func(s string) bool {
		_, err := strconv.ParseFloat(s, 64)
		return err == nil
	}

	isInteger := func(s string) bool {
		_, err := strconv.ParseInt(s, 10, 64)
		return err == nil
	}
	if len(ts.word) > 0 {
		w := string(ts.word)
		k := ts.kind
		if k == AString {
			// Don't convert if actual string.
		} else if isInteger(w) {
			k = AInteger
		} else if isFloat(w) {
			k = AFloat
		}
		ts.Tokens = append(ts.Tokens, Token{k, w})
		ts.word = make([]rune, 0)
		ts.kind = ASymbol
	}
}

func (ts *Tokens) setKind(kind tokenType) {
	ts.kind = kind
}

func (ts *Tokens) pushToken(kind tokenType, value string) {
	ts.Tokens = append(ts.Tokens, Token{kind, value})
}

func (ts *Tokens) inString() bool {
	return ts.kind == AString
}

func (ts *Tokens) emptyWord() bool {
	return len(ts.word) == 0
}

// Tokenize a line of code.
func Tokenize(form string) (*Tokens, error) {
	results := &Tokens{
		Tokens: make([]Token, 0),
		word:   make([]rune, 0),
		form:   form,
		kind:   ASymbol,
	}

	for _, c := range form {
		switch c {

		case '(':
			results.pushToken(AOpenParen, "(")

		case ')':
			results.pushWord()
			results.pushToken(ACloseParen, ")")

		case '"':
			if results.emptyWord() {
				results.setKind(AString)
			} else {
				results.pushWord()
			}

		case ',': // Treat commas as whitespace.
			if results.inString() {
				results.pushChar(c)
			}

		case ' ', '\t', '\r', '\n':
			if results.inString() {
				results.pushChar(c)
			} else {
				results.pushWord()
			}

		case '\'':
			if results.emptyWord() {
				results.pushToken(AQuote, "'")
			} else {
				results.pushChar(c)
			}

		default:
			results.pushChar(c)
		}
	}
	results.pushWord()

	return results, nil
}

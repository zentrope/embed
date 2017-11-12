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

package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

// Sexp represents a parse s-expression
type Sexp interface {
	AsString() string
}

type sexpPrimitive primitiveFunc
type sexpList []Sexp
type sexpString string
type sexpInteger int64
type sexpFloat float64
type sexpSymbol string
type sexpQuote struct {
	quote Sexp
}

func (s sexpPrimitive) AsString() string {
	return fmt.Sprintf("builtin::%v", s)
}

func (s sexpQuote) AsString() string {
	return "(quote " + s.quote.AsString() + ")"
}

func (s sexpList) AsString() string {
	elems := make([]string, 0)
	for _, e := range s {
		elems = append(elems, e.AsString())
	}
	return fmt.Sprintf("(%v)", strings.Join(elems, " "))
}

func (s sexpString) AsString() string {
	return "str::" + string(s)
}

func (s sexpInteger) AsString() string {
	return fmt.Sprintf("int::%d", s)
}

func (s sexpFloat) AsString() string {
	return fmt.Sprintf("float::%f", s)
}

func (s sexpSymbol) AsString() string {
	return "sym::" + string(s)
}

// Parser represents the state of the parser.
type Parser struct {
	tokens   []Token
	position int
}

// NewParser returns a new instance of a Parser.
func NewParser(tokens *Tokens) *Parser {
	return &Parser{
		tokens:   tokens.Tokens,
		position: 0,
	}
}

func (p *Parser) pushBack() {
	if p.position > 0 {
		p.position = p.position - 1
	}
}

func (p *Parser) next() Token {
	t := p.tokens[p.position]
	p.position = p.position + 1
	return t
}

func (p *Parser) notDone() bool {
	return p.position+1 != len(p.tokens)
}

// Parse returns an s-expression suitable for interpretation.
func (p *Parser) Parse() (Sexp, error) {
	token := p.next()

	switch token.kind {

	case AOpenParen:
		return p.parseList()

	case ASymbol:
		return sexpSymbol(token.value), nil

	case AString:
		return sexpString(token.value), nil

	case AInteger:
		i, _ := strconv.ParseInt(token.value, 10, 64)
		return sexpInteger(i), nil

	case AFloat:
		f, _ := strconv.ParseFloat(token.value, 64)
		return sexpFloat(f), nil

	case AQuote:
		sexp, err := p.Parse()
		if err != nil {
			return sexp, err
		}
		return sexpQuote{sexp}, nil

	default:
		return sexpString("err"),
			fmt.Errorf("unable to process token '%v'", token)
	}
}

func (p *Parser) parseList() (Sexp, error) {
	list := make([]Sexp, 0)

	for p.notDone() {
		token := p.next()

		switch token.kind {

		case AOpenParen:
			sublist, err := p.parseList()
			if err != nil {
				return sublist, err
			}
			list = append(list, sublist)

		case ACloseParen:
			break

		default:
			p.pushBack()
			atom, err := p.Parse()
			if err != nil {
				return atom, err
			}
			list = append(list, atom)
		}
	}
	return sexpList(list), nil
}

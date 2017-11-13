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
)

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

// Parse2 returns an s-expression suitable for interpretation.
func (p *Parser) Parse() (Expression, error) {
	token := p.next()
	switch token.kind {

	case AOpenParen:
		return p.parseList()

	case ASymbol:
		return NewExpr(ExpSymbol, token.value), nil

	case AString:
		return NewExpr(ExpString, token.value), nil

	case AInteger:
		i, _ := strconv.ParseInt(token.value, 10, 64)
		return NewExpr(ExpInteger, i), nil

	case AFloat:
		f, _ := strconv.ParseFloat(token.value, 64)
		return NewExpr(ExpFloat, f), nil

	case AQuote:
		sexp, err := p.Parse()
		if err != nil {
			return sexp, err
		}
		return NewExpr(ExpQuote, sexp), nil

	default:
		return NilExpression, fmt.Errorf("unable to process token '%v'", token)
	}
}

func (p *Parser) parseList() (Expression, error) {
	list := make([]Expression, 0)

done:
	for p.notDone() {
		token := p.next()

		switch token.kind {

		case AOpenParen:
			sublist, err := p.parseList()
			if err != nil {
				return NilExpression, err
			}
			list = append(list, sublist)

		case ACloseParen:
			break done

		default:
			p.pushBack()
			atom, err := p.Parse()
			if err != nil {
				return atom, err
			}
			list = append(list, atom)
		}
	}
	return NewExpr(ExpList, list), nil
}

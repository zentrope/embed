//
// Copyright © 2017-present Keith Irwin
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

import "fmt"

// Interpreter is something that evaluates
type Interpreter interface {
	Evaluate(env *Environment, expr Expression) (Expression, error)
	Execute(form string) (Expression, error)
}

// TcoInterpreter attempts to implement some TCO
type TcoInterpreter struct {
	parser      *Parser
	environment *Environment
}

// NaiveInterpreter is fully recursive
type NaiveInterpreter struct {
	parser      *Parser
	environment *Environment
}

// Type represents a type of evaluator
type Type int

// EvaluatorTypes
const (
	TCO Type = iota
	Naive
)

// NewInterpreter returns an evaluator
func NewInterpreter(kind Type) Interpreter {
	switch kind {
	case TCO:
		return TcoInterpreter{
			environment: NewEnvironment(),
			parser:      NewParser(),
		}
	default:
		return NaiveInterpreter{
			environment: NewEnvironment(),
			parser:      NewParser(),
		}
	}
}

// Execute a Haki expression.
func (tco TcoInterpreter) Execute(form string) (Expression, error) {

	tokens, err := Tokenize(form)
	if err != nil {
		return NilExpression, err
	}

	tco.parser.Reset(tokens)

	expr, err := tco.parser.Parse()
	if err != nil {
		return NilExpression, err
	}

	return tco.Evaluate(tco.environment, expr)
}

// Execute a Haki expression.
func (naive NaiveInterpreter) Execute(form string) (Expression, error) {

	tokens, err := Tokenize(form)
	if err != nil {
		return NilExpression, err
	}

	naive.parser.Reset(tokens)

	expr, err := naive.parser.Parse()
	if err != nil {
		return NilExpression, err
	}
	return naive.Evaluate(naive.environment, expr)
}

func nilExpr(reason string, params ...interface{}) (Expression, error) {
	return NilExpression, fmt.Errorf(reason, params...)
}

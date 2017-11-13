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

package interpreter

import (
	"fmt"
)

func apply(env Environment, op Expression, args []Expression) (Expression, error) {

	theOp, err := Evaluate(env, op)
	if err != nil {
		return NilExpression, err
	}

	params := make([]Expression, 0)
	for _, a := range args {
		param, err := Evaluate(env, a)
		if err != nil {
			return NilExpression, err
		}
		params = append(params, param)
	}

	if theOp.IsPrimitive() {
		result, err := theOp.InvokePrimitive(params)
		if err != nil {
			return NilExpression, err
		}
		return result, nil
	}

	return NilExpression, fmt.Errorf("function not found: '%v'", theOp)
}

func evalIf(env Environment, exprs Expression) (Expression, error) {
	argc := len(exprs.list)
	if argc < 2 {
		return NilExpression, fmt.Errorf("too few arguments (%v) to if", argc)
	}
	if argc > 3 {
		return NilExpression, fmt.Errorf("too many arguments (%v) to if", argc)
	}

	test, err := Evaluate(env, exprs.list[0])
	if err != nil {
		return NilExpression, err
	}

	var result Expression

	if test.IsTruthy() {
		result, err = Evaluate(env, exprs.list[1])
		if err != nil {
			return NilExpression, err
		}
	} else if argc >= 3 {
		result, err = Evaluate(env, exprs.list[2])
		if err != nil {
			return NilExpression, err
		}
	}
	return result, nil
}

func evalAnd(env Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if !result.IsTruthy() {
			return result, nil
		}
	}

	return result, nil
}

func evalOr(env Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if result.IsTruthy() {
			return result, nil
		}
	}
	return result, nil
}

// Evaluate an expression
func Evaluate(env Environment, expr Expression) (Expression, error) {

	if expr.IsSymbol() {
		found, value := env.Lookup(expr.symbol)
		if !found {
			return NilExpression, fmt.Errorf("value not found for '%v'", expr.String())
		}
		return Evaluate(env, value)
	}

	if expr.IsQuote() {
		return *expr.quote, nil
	}

	if expr.IsAtom() {
		return expr, nil
	}

	if expr.IsList() {
		if expr.StartsWith("and") {
			return evalAnd(env, expr.Tail())
		}
		if expr.StartsWith("or") {
			return evalOr(env, expr.Tail())
		}
		if expr.StartsWith("if") {
			return evalIf(env, expr.Tail())
		}
		return apply(env, expr.Head(), expr.Tail().list)
	}

	return NilExpression, fmt.Errorf("unable to eval expression [%v]", expr)
}

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

func apply(env *Environment, op Sexp, args []Sexp) (Sexp, error) {

	theOp, err := Evaluate(env, op)
	if err != nil {
		return nil, err
	}

	params := make([]Sexp, 0)
	for _, a := range args {
		param, err := Evaluate(env, a)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}

	switch f := theOp.(type) {
	case sexpPrimitive:
		result, err := f(params)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, fmt.Errorf("function not found: '%v'", f)
	}
}

func isSpecial(key string, expr Sexp) bool {
	switch t := expr.(type) {
	case sexpSymbol:
		return key == string(t)
	default:
		return false
	}
}

func isTruthy(sexp Sexp) bool {
	if sexp == nil {
		return false
	}

	switch t := sexp.(type) {
	case sexpBool:
		return bool(t)
	}
	return true
}

func evalIf(env *Environment, exprs []Sexp) (Sexp, error) {
	if len(exprs) < 2 {
		return nil, fmt.Errorf("too few arguments (%v) to if", len(exprs))
	}
	if len(exprs) > 3 {
		return nil, fmt.Errorf("too many arguments (%v) to if", len(exprs))
	}

	test, err := Evaluate(env, exprs[0])
	if err != nil {
		return nil, err
	}

	var result Sexp

	if isTruthy(test) {
		result, err = Evaluate(env, exprs[1])
		if err != nil {
			return nil, err
		}
	} else if len(exprs) >= 3 {
		result, err = Evaluate(env, exprs[2])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalAnd(env *Environment, exprs []Sexp) (Sexp, error) {

	var result Sexp
	var err error

	for _, e := range exprs {
		result, err = Evaluate(env, e)
		if err != nil {
			return nil, err
		}

		if !isTruthy(result) {
			return result, nil
		}
	}

	return result, nil
}

func evalOr(env *Environment, exprs []Sexp) (Sexp, error) {

	var result Sexp
	var err error

	for _, e := range exprs {
		result, err = Evaluate(env, e)
		if err != nil {
			return nil, err
		}

		if isTruthy(result) {
			return result, nil
		}
	}
	return result, nil
}

func head(list []Sexp) Sexp {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func tail(list []Sexp) []Sexp {
	if len(list) == 0 {
		return list
	}
	return list[1:]
}

// Evaluate an expression
func Evaluate(env *Environment, expr Sexp) (Sexp, error) {

	switch t := expr.(type) {

	case sexpList:
		op := head(t)
		if isSpecial("and", op) {
			return evalAnd(env, tail(t))
		}
		if isSpecial("or", op) {
			return evalOr(env, tail(t))
		}
		if isSpecial("if", op) {
			return evalIf(env, tail(t))
		}
		return apply(env, head(t), tail(t))

	case sexpInteger:
		return t, nil

	case sexpFloat:
		return t, nil

	case sexpString:
		return t, nil

	case sexpQuote:
		return t.quote, nil

	case sexpSymbol:

		found, value := env.Lookup(string(t))
		if !found {
			return nil, fmt.Errorf("value not found for '%v'", t.AsString())
		}
		return Evaluate(env, value)

	case sexpPrimitive:
		return t, nil

	default:
		return nil, fmt.Errorf("can't parse [%v]", t)
	}
}

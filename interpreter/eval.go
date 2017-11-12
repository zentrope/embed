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

	fmt.Printf(" op: %#v\n", theOp)
	fmt.Printf(" args: %#v\n", params)

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

	//return nil, errors.New("not implemented")
}

func first(list []Sexp) Sexp {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func rest(list []Sexp) []Sexp {
	if len(list) == 0 {
		return list
	}
	return list[1:]
}

// Evaluate an expression
func Evaluate(env *Environment, expr Sexp) (Sexp, error) {
	switch t := expr.(type) {

	case sexpList:
		// apply first to rest
		return apply(env, first(t), rest(t))

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
			return nil, fmt.Errorf("value not found for '%v' (%#v)", t, t)
		}
		return Evaluate(env, value)

	case sexpPrimitive:
		return t, nil

	default:
		return nil, fmt.Errorf("can't parse [%v]", t)
	}
}

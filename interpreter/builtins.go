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
	"errors"
	"fmt"
	"strings"
)

var builtins = map[string]primitiveFunc{
	"+":   primitiveAdd,
	"-":   primitiveMinus,
	"prn": primitivePrn,
	"=":   pEquals,
}

type primitiveFunc func(args []Sexp) (Sexp, error)

// TRUE represents a true boolean value
const TRUE = sexpBool(true)

// FALSE represents a false boolean value
const FALSE = sexpBool(false)

func pEquals(args []Sexp) (Sexp, error) {
	// Return true if all the arguments are equal to each other in value
	// and type.

	if len(args) < 1 {
		return FALSE, errors.New("wrong number of args for '=', must be at least one")
	}

	sentinel := args[0]

	for _, a := range args[1:] {
		if a != sentinel {
			return FALSE, nil
		}
	}
	return TRUE, nil
}

func primitiveAdd(args []Sexp) (Sexp, error) {
	var result float64
	for _, arg := range args {
		switch x := arg.(type) {
		case sexpFloat:
			result = result + float64(x)
		case sexpInteger:
			result = result + float64(x)
		default:
			return nil, fmt.Errorf("unknown argument type, int/float expected, got [%#v]", x)
		}
	}
	return sexpFloat(result), nil
}

func primitiveMinus(args []Sexp) (Sexp, error) {
	var result float64
	for _, arg := range args {
		switch x := arg.(type) {
		case sexpFloat:
			result = result - float64(x)
		case sexpInteger:
			result = result - float64(x)
		default:
			return nil, fmt.Errorf("unknown argument type, int/float expected, got [%#v]", x)
		}
	}
	return sexpFloat(result), nil
}

func primitivePrn(args []Sexp) (Sexp, error) {
	values := make([]string, 0)
	for _, a := range args {
		values = append(values, fmt.Sprintf("%v", a))
	}
	fmt.Println(strings.Join(values, " "))
	return nil, nil
}

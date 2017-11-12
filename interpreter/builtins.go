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
	"strings"
)

var builtins = map[string]primitiveFunc{
	"+":   primitiveAdd,
	"-":   primitiveMinus,
	"prn": primitivePrn,
}

type primitiveFunc func(args []Sexp) (Sexp, error)

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

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

import (
	"errors"
	"fmt"
	"strings"
)

type primitivesMap map[string]primitiveFunc

var defaultBuiltins = primitivesMap{
	"=":   _equals,
	"not": _not,
	"prn": _prn,
}

var builtins = make(primitivesMap, 0)

func init() {
	prims := []primitivesMap{
		defaultBuiltins,
		mathBuiltins,    // builtins_math
		stringBuiltins,  // builtins_string
		listBuiltins,    // builtins_list
		fileioBuiltins,  // builtins_fileio
		hashmapBuiltins, // builtins_hashmap
	}
	for _, prim := range prims {
		for name, fn := range prim {
			builtins[name] = fn
		}
	}
}

type primitiveFunc func(args []Expression) (Expression, error)

func isIntegral(val float64) bool {
	return val == float64(int64(val))
}

func verifyNums(args []Expression) error {
	for _, arg := range args {
		switch arg.tag {
		case ExpInteger, ExpFloat:
			continue
		default:
			return errors.New("all arguments must be numbers")
		}
	}
	return nil
}

func verifyStrings(args []Expression) error {
	for _, arg := range args {
		switch arg.tag {
		case ExpString:
			continue
		default:
			return errors.New("all arguments must be strings")
		}
	}
	return nil
}

func asNumber(expr Expression) (float64, error) {
	switch expr.tag {
	case ExpInteger:
		return float64(expr.integer), nil
	case ExpFloat:
		return expr.float, nil
	default:
		return 0, errors.New("not a number")
	}
}

func toExpr(x float64) Expression {
	if isIntegral(x) {
		return NewExpr(ExpInteger, int64(x))
	}
	return NewExpr(ExpFloat, x)
}

func _not(args []Expression) (Expression, error) {
	if len(args) > 1 {
		return nilExpr("(not expr) takes one parameter, you provided %v", len(args))
	}

	return NewExpr(ExpBool, !args[0].IsTruthy()), nil
}

func _equals(args []Expression) (Expression, error) {
	// Return true if all the arguments are equal to each other in value
	// and type.

	if len(args) < 1 {
		return FalseExpression, errors.New("wrong number of args for '=', must be at least one")
	}

	sentinel := args[0]

	for _, a := range args[1:] {
		if !a.Equals(sentinel) {
			return FalseExpression, nil
		}
	}
	return TrueExpression, nil
}

func _prn(args []Expression) (Expression, error) {
	values := make([]string, 0)
	for _, a := range args {
		value := a.String()
		if a.tag == ExpString {
			value = a.string
		}
		values = append(values, value)
	}
	fmt.Println(strings.Join(values, " "))
	return NilExpression, nil
}

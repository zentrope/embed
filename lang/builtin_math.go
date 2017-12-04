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
	"math"
)

var mathBuiltins = primitivesMap{
	"*":   _mult,
	"+":   _add,
	"-":   _minus,
	"<":   _lessThan,
	"mod": _mod,
}

func isIntegral(val float64) bool {
	return val == float64(int64(val))
}

func toExpr(x float64) Expression {
	if isIntegral(x) {
		return NewExpr(ExpInteger, int64(x))
	}
	return NewExpr(ExpFloat, x)
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

func _mult(args []Expression) (Expression, error) {
	if err := verifyNums(args); err != nil {
		return NilExpression, err
	}
	result := 1.0
	for _, arg := range args {
		num, err := asNumber(arg)
		if err != nil {
			return NilExpression, err
		}
		result = result * num
	}
	return toExpr(result), nil
}

func _lessThan(args []Expression) (Expression, error) {

	argc := len(args)

	if argc < 1 {
		return nilExpr("(< a b ... n) requires at least 1 arg")
	}

	if err := verifyNums(args); err != nil {
		return NilExpression, err
	}

	sentinel, err := asNumber(args[0])
	if err != nil {
		return NilExpression, err
	}

	for _, arg := range args[1:] {
		candidate, err := asNumber(arg)
		if err != nil {
			return NilExpression, err
		}
		if candidate <= sentinel {
			return FalseExpression, nil
		}
		sentinel = candidate
	}

	return TrueExpression, nil
}

func _mod(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 2 {
		return nilExpr("(mod num div) takes 2 args, you provided %v", argc)
	}

	num := args[0]
	div := args[1]

	n := num.float
	d := div.float

	if num.tag == ExpInteger {
		n = float64(num.integer)
	}
	if div.tag == ExpInteger {
		d = float64(div.integer)
	}

	result := math.Mod(n, d)

	if isIntegral(result) {
		return NewExpr(ExpInteger, int64(result)), nil
	}
	return NewExpr(ExpFloat, result), nil
}

func _add(args []Expression) (Expression, error) {
	var result float64
	for _, arg := range args {
		switch arg.tag {
		case ExpFloat:
			result = result + float64(arg.float)
		case ExpInteger:
			result = result + float64(arg.integer)
		default:
			return nilExpr("unknown argument type for [%v], [int/float] expected, got [%v]", arg, arg.Type())
		}
	}

	if isIntegral(result) {
		return NewExpr(ExpInteger, int64(result)), nil
	}
	return NewExpr(ExpFloat, result), nil
}

func _minus(args []Expression) (Expression, error) {

	if len(args) < 1 {
		return NilExpression, errors.New("`-` requires 1 or more args")
	}

	var result float64

	switch args[0].tag {
	case ExpFloat:
		result = float64(args[0].float)
	case ExpInteger:
		result = float64(args[0].integer)
	default:
		return nilExpr("In '-', unknown argument type [%v], [int/float] expected, got [%v]", args[0], args[0].Type())
	}

	if len(args) == 1 {
		result = -1.0 * result

		if isIntegral(result) {
			return NewExpr(ExpInteger, int64(result)), nil
		}
		return NewExpr(ExpFloat, result), nil
	}

	for _, arg := range args[1:] {
		switch arg.tag {
		case ExpFloat:
			result = result - float64(arg.float)
		case ExpInteger:
			result = result - float64(arg.integer)
		default:
			return nilExpr("In '-', unknown argument type [%v], [int/float] expected, got [%v]", arg, arg.Type())
		}
	}

	if isIntegral(result) {
		return NewExpr(ExpInteger, int64(result)), nil
	}
	return NewExpr(ExpFloat, result), nil
}

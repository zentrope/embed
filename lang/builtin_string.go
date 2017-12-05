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
	"fmt"
	"regexp"
	"strings"
)

var stringBuiltins = primitivesMap{
	"format":   _format,
	"re-find":  _reFind,
	"re-list":  _reList,
	"re-match": _reMatch,
	"re-split": _reSplit,
	"trim":     _trim,
	"triml":    _triml,
	"trimr":    _trimr,
}

var _trimCutSet = " \n\t\r"

func _trim(args []Expression) (Expression, error) {

	if err := typeCheck("(trim string)", args, ckArity(1), ckString(0)); err != nil {
		return NilExpression, err
	}

	return NewStringExpr(strings.TrimSpace(args[0].string)), nil
}

func _triml(args []Expression) (Expression, error) {

	if err := typeCheck("(triml string)", args, ckArity(1), ckString(0)); err != nil {
		return NilExpression, err
	}

	return NewStringExpr(strings.TrimLeft(args[0].string, _trimCutSet)), nil
}

func _trimr(args []Expression) (Expression, error) {

	if err := typeCheck("(trimr string)", args, ckArity(1), ckString(0)); err != nil {
		return NilExpression, err
	}

	return NewStringExpr(strings.TrimRight(args[0].string, _trimCutSet)), nil
}

func _format(args []Expression) (Expression, error) {
	sig := "(format pattern v ... vs)"
	specs := []spec{ckArityAtLeast(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	pattern := args[0].string
	params := make([]interface{}, 0)
	for _, a := range args[1:] {
		params = append(params, a.Value())
	}

	result := fmt.Sprintf(pattern, params...)
	return NewExpr(ExpString, result), nil
}

func _reFind(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 2 {
		return nilExpr("(re-find re string) requires 2 args, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
		return NilExpression, err
	}

	re, err := regexp.Compile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpString, re.FindString(args[1].string)), nil
}

func _reMatch(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 2 {
		return nilExpr("(re-match re string) requires 2 args, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
		return NilExpression, err
	}

	re, err := regexp.Compile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpBool, re.MatchString(args[1].string)), nil
}

func _reSplit(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 2 {
		return nilExpr("(re-split re string) requires 2 args, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
		return NilExpression, err
	}

	re, err := regexp.Compile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	words := re.Split(args[1].string, -1)

	es := make([]Expression, 0)

	for _, word := range words {
		es = append(es, NewExpr(ExpString, word))
	}

	return NewExpr(ExpList, es), nil
}

func _reList(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 2 {
		return nilExpr("(re-list re string) requires 2 args, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
		return NilExpression, err
	}

	re, err := regexp.Compile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	words := re.FindAllString(args[1].string, -1)

	es := make([]Expression, 0)

	for _, word := range words {
		es = append(es, NewExpr(ExpString, word))
	}

	return NewExpr(ExpList, es), nil
}

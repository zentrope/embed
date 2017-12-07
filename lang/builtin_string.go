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
	"ends-with?":   _endsWithP,
	"index":        _index,
	"format":       _format,
	"last-index":   _lastIndex,
	"lower-case":   _lowerCase,
	"re-find":      _reFind,
	"re-list":      _reList,
	"re-match":     _reMatch,
	"re-split":     _reSplit,
	"replace":      _replace,
	"starts-with?": _startsWithP,
	"substr":       _substr,
	"trim":         _trim,
	"triml":        _triml,
	"trimr":        _trimr,
	"upper-case":   _upperCase,
}

// CoreStringFunctions written in Haki
var CoreStringFunctions = `
	(defun words (s)
		(re-split "\s+" s))
`

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

func _lowerCase(args []Expression) (Expression, error) {
	if err := typeCheck("(lower-case s)", args, ckArity(1), ckString(0)); err != nil {
		return NilExpression, err
	}

	return NewStringExpr(strings.ToLower(args[0].string)), nil
}

func _replace(args []Expression) (Expression, error) {
	if err := typeCheck("(replace s old new)", args, ckArity(3), ckString(0, 1, 2)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	old := args[1].string
	new := args[2].string

	return NewStringExpr(strings.Replace(s, old, new, -1)), nil
}

func _substr(args []Expression) (Expression, error) {
	if err := typeCheck("(substr s start end)", args, ckArity(3), ckString(0), ckInt(1, 2)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	start := args[1].integer
	end := args[2].integer

	if start < 0 {
		return NilExpression, fmt.Errorf("(substr s start end) → `start` (%v) should be a positive int",
			start)
	}

	if (end - start) < 0 {
		return NilExpression, fmt.Errorf("(substr s start end) → `end` (%v) should be >= to `start` (%v)",
			end, start)
	}

	max := int64(len(s))
	if end > max {
		return NilExpression, fmt.Errorf("(substr s start end) → provided value for `end` (%v) exceeds `s` length `%v`",
			end, len(s))
	}

	if start < 0 {
		return NilExpression, fmt.Errorf("(substr s start end) → provided value for `end` (%v) exceeds `s` length `%v`",
			end, len(s))
	}

	cut := s[start:end]

	return NewStringExpr(cut), nil
}

func _startsWithP(args []Expression) (Expression, error) {
	if err := typeCheck("(starts-with? s)", args,
		ckArity(2), ckString(0, 1)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	prefix := args[1].string

	return NewBoolExpr(strings.HasPrefix(s, prefix)), nil
}

func _endsWithP(args []Expression) (Expression, error) {
	if err := typeCheck("(ends-with? s)", args,
		ckArity(2), ckString(0), ckString(1)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	suffix := args[1].string

	return NewBoolExpr(strings.HasSuffix(s, suffix)), nil
}

func _upperCase(args []Expression) (Expression, error) {
	if err := typeCheck("(upper-case s)", args, ckArity(1), ckString(0)); err != nil {
		return NilExpression, err
	}

	return NewStringExpr(strings.ToUpper(args[0].string)), nil
}

func _index(args []Expression) (Expression, error) {
	if err := typeCheck("(index s substr)", args, ckArity(2), ckString(0, 1)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	substr := args[1].string

	val := strings.Index(s, substr)
	return NewIntExpr(int64(val)), nil
}

func _lastIndex(args []Expression) (Expression, error) {
	if err := typeCheck("(last-index s substr)", args, ckArity(2), ckString(0, 1)); err != nil {
		return NilExpression, err
	}

	s := args[0].string
	substr := args[1].string

	val := strings.LastIndex(s, substr)
	return NewIntExpr(int64(val)), nil
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

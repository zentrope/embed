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
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strings"
)

type primitivesMap map[string]primitiveFunc

var defaultBuiltins = map[string]primitiveFunc{
	"*":          _mult,
	"+":          _add,
	"-":          _minus,
	"<":          _lessThan,
	"=":          _equals,
	"append":     _append,
	"close":      _close,
	"count":      _count,
	"directory?": _directoryP,
	"exists?":    _existsP,
	"file?":      _fileP,
	"format":     _format,
	"handle?":    _handleP,
	"head":       _head,
	"join":       _join,
	"list":       _list,
	"list?":      _listP,
	"read-file":  _readFile,
	"mod":        _mod,
	"not":        _not,
	"open":       _open,
	"prepend":    _prepend,
	"prn":        _prn,
	"re-find":    _reFind,
	"re-list":    _reList,
	"re-match":   _reMatch,
	"re-split":   _reSplit,
	"tail":       _tail,
}

var builtins = make(primitivesMap, 0)

func init() {
	prims := []primitivesMap{
		defaultBuiltins,
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

func _format(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 1 {
		return nilExpr("(format pattern ...args) requires at least the first argument.")
	}

	if args[0].tag != ExpString {
		return nilExpr("(format pattern ...args) the pattern argument must be a string.")
	}

	pattern := args[0].string
	params := make([]interface{}, 0)
	for _, a := range args[1:] {
		params = append(params, a.Value())
	}

	result := fmt.Sprintf(pattern, params...)
	return NewExpr(ExpString, result), nil
}

//-----------------------------------------------------------------------------
// FILE BUILTINS
//-----------------------------------------------------------------------------

func _readFile(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 1 {
		return nilExpr("(read-file file-name) requires 1 arg, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
		return NilExpression, err
	}

	buffer, err := ioutil.ReadFile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	str := string(buffer)
	return NewExpr(ExpString, str), nil
}

func _open(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(open file-path) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(open file-path) expects file-path to be a string, not %v",
			ExprTypeName(args[0].tag))
	}
	path := args[0].string
	file, err := os.Open(path)
	if err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpFile, file), nil
}

func _close(args []Expression) (Expression, error) {

	argc := len(args)
	if argc != 1 {
		return nilExpr("(close file-handle) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpFile {
		return nilExpr("(close file-handle) expects file-handle to be a file-handle, not %v",
			ExprTypeName(args[0].tag))
	}

	fileData := args[0].file
	fileData.isOpen = false
	if err := fileData.file.Close(); err != nil {
		return NilExpression, err
	}

	return NilExpression, nil
}

func _directoryP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(directory? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(directory? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
	}

	path := args[0].string

	f, err := os.Open(path)
	if err != nil {
		return FalseExpression, nil
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FalseExpression, nil
	}

	return NewExpr(ExpBool, info.IsDir()), nil
}

func _existsP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(exists? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(exists? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
	}

	path := args[0].string

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return TrueExpression, nil
	}
	return FalseExpression, nil
}

func _fileP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(file? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(file? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
	}

	path := args[0].string

	f, err := os.Open(path)
	if err != nil {
		return FalseExpression, nil
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FalseExpression, nil
	}

	return NewExpr(ExpBool, !info.IsDir()), nil
}

func _handleP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(handle? file-handle) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpFile {
		return nilExpr("(handle? file-handle) ← file-handle should be 'handle', not '%v'",
			ExprTypeName(args[0].tag))
	}

	return TrueExpression, nil
}

//-----------------------------------------------------------------------------
// STRING / REGEX BUILTINS
//-----------------------------------------------------------------------------

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

func _count(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 1 {
		return nilExpr("(count lst) ← lst should be a list")
	}

	return NewExpr(ExpInteger, int64(len(args[0].list))), nil
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

func _not(args []Expression) (Expression, error) {
	if len(args) > 1 {
		return nilExpr("(not expr) takes one parameter, you provided %v", len(args))
	}

	return NewExpr(ExpBool, !args[0].IsTruthy()), nil
}

func _list(args []Expression) (Expression, error) {
	return NewExpr(ExpList, args), nil
}

func _listP(args []Expression) (Expression, error) {
	if len(args) != 1 {
		return nilExpr("(list? val) requires one argument, you provided: %v.", len(args))
	}

	if args[0].tag == ExpList {
		return TrueExpression, nil
	}
	return FalseExpression, nil
}

// (prepend x list)
func _prepend(args []Expression) (Expression, error) {

	if len(args) < 2 {
		return nilExpr("(prepend val list) requires two params: item, list")
	}

	item := args[0]
	list := args[1]

	if list.tag != ExpList {
		return nilExpr("(prepend val list) 2nd parameter must be a 'list', not '%v'",
			ExprTypeName(list.tag),
		)
	}

	return NewExpr(ExpList, append([]Expression{item}, list.list...)), nil
}

// (append list x)
func _append(args []Expression) (Expression, error) {
	if len(args) != 2 {
		return nilExpr("append takes two args (list, item), you provided %v", len(args))
	}

	list := args[0]
	item := args[1]

	if list.tag != ExpList {
		return nilExpr("append's first arg (list, item) must be a list")
	}

	return NewExpr(ExpList, append(list.list, item)), nil
}

// (join list1 list2 ... listn)
func _join(args []Expression) (Expression, error) {

	for _, e := range args {
		if e.tag != ExpList {
			return nilExpr("join takes only list params, %v is not a list", e)
		}
	}

	newList := make([]Expression, 0)
	for _, l := range args {
		newList = append(newList, l.list...)
	}

	return NewExpr(ExpList, newList), nil
}

func _head(args []Expression) (Expression, error) {
	if len(args) == 0 {
		return NilExpression, errors.New("head requires a parameter")
	}

	if !args[0].IsList() {
		return NilExpression, errors.New("head requires a list parameter")
	}

	list := args[0].list
	if len(list) == 0 {
		return NilExpression, nil
	}
	return list[0], nil
}

func _tail(args []Expression) (Expression, error) {

	if len(args) == 0 {
		return NilExpression, errors.New("tail requires a parameter")
	}

	if !args[0].IsList() {
		return NilExpression, errors.New("tail requires a list parameter")
	}

	list := args[0].list

	if len(list) == 0 {
		return NilExpression, nil
	}
	return NewExpr(ExpList, list[1:]), nil
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

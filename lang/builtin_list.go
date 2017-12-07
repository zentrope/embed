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
)

var listBuiltins = map[string]primitiveFunc{
	"append":  _append,
	"count":   _count,
	"head":    _head,
	"join":    _join,
	"list":    _list,
	"list?":   _listP,
	"prepend": _prepend,
	"tail":    _tail,
}

// NewListExpr constructs a new list
func NewListExpr(list []Expression) Expression {

	data := make([]interface{}, 0)
	data = append(data, ExpList)

	for _, e := range list {
		data = append(data, e.hash)
	}

	return Expression{tag: ExpList, hash: hashIt(data...), list: list}
}

//-----------------------------------------------------------------------------
// IMPLEMENTATION
//-----------------------------------------------------------------------------

func _count(args []Expression) (Expression, error) {

	if err := typeCheck("(count string|list|hash-map)", args,
		ckArity(1), ckCountable(0)); err != nil {
		return NilExpression, err
	}

	e := args[0]
	var c int

	if e.IsList() {
		c = len(e.list)
	} else if e.IsHashMap() {
		c = len(e.hashMap.keys)
	} else {
		c = len(e.string)
	}

	return NewIntExpr(int64(c)), nil
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

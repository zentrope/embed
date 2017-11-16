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
	"hash/fnv"
	"log"
	"strings"
)

// ExpressionType is the type of expression
type ExpressionType int

// ExpressionTypes
const (
	ExpNil ExpressionType = iota
	ExpPrimitive
	ExpFunction
	ExpList
	ExpString
	ExpInteger
	ExpFloat
	ExpSymbol
	ExpBool
	ExpQuote
)

// Expression represents a computation
type Expression struct {
	tag            ExpressionType
	hash           uint32
	string         string
	integer        int64
	float          float64
	symbol         string
	bool           bool
	list           []Expression
	quote          *Expression
	primitive      primitiveFunc
	functionName   string
	functionParams *Expression
	functionBody   *Expression
}

func hashIt(values ...interface{}) uint32 {
	hash := fnv.New32()
	for _, v := range values {
		fmt.Fprint(hash, v)
	}
	return hash.Sum32()
}

// WrapImplicitDo wraps an expression in a do block
func WrapImplicitDo(body []Expression) Expression {
	do := NewExpr(ExpSymbol, "do")
	return NewExpr(ExpList, append([]Expression{do}, body...))
}

// NewFunctionExpr returns an expression representing a function
func NewFunctionExpr(name Expression, params Expression, body Expression) Expression {
	p := params
	b := WrapImplicitDo(body.list)

	return Expression{
		tag:            ExpFunction,
		hash:           hashIt(ExpFunction, name, p.hash, b.hash),
		functionName:   name.symbol,
		functionParams: &p,
		functionBody:   &b,
	}
}

// NewExpr constructs a new expression of the given type
func NewExpr(tag ExpressionType, value interface{}) Expression {

	data := make([]interface{}, 0)
	data = append(data, tag)

	if tag == ExpList {
		for _, e := range value.([]Expression) {
			data = append(data, e.hash)
		}
	} else {
		data = append(data, value)
	}

	e := Expression{tag: tag, hash: hashIt(data...)}
	switch tag {
	case ExpPrimitive:
		e.primitive = value.(primitiveFunc)
	case ExpList:
		e.list = value.([]Expression)
	case ExpString:
		e.string = value.(string)
	case ExpInteger:
		e.integer = value.(int64)
	case ExpFloat:
		e.float = value.(float64)
	case ExpSymbol:
		e.symbol = value.(string)
	case ExpBool:
		e.bool = value.(bool)
	case ExpQuote:
		exp := value.(Expression)
		e.quote = &exp
	default:
		log.Fatalf("unable to create new expr of type %v", tag)
	}
	return e
}

// NilExpression represents nil
var NilExpression = Expression{tag: ExpNil}

// TrueExpression for a boolean true
var TrueExpression = NewExpr(ExpBool, true)

// FalseExpression for a boolean false
var FalseExpression = NewExpr(ExpBool, false)

// StartsWith returns true if first elem in list is named prefix.
func (e Expression) StartsWith(prefix string) bool {
	if e.tag != ExpList {
		return false
	}

	elem := e.Head()
	return prefix == elem.symbol
}

// Head returns the first element of the list.
func (e Expression) Head() Expression {
	if e.tag != ExpList {
		log.Fatalf("Can't take the head of a %v", e)
	}

	if len(e.list) == 0 {
		return NilExpression
	}
	return e.list[0]
}

// Tail returns the rest of the elements of a list.
func (e Expression) Tail() Expression {
	if e.tag != ExpList {
		log.Fatalf("Can't take the tail of a %v", e)
	}

	if len(e.list) == 0 {
		return e
	}
	return Expression{tag: ExpList, list: e.list[1:]}
}

func (e Expression) String() string {
	switch e.tag {
	case ExpPrimitive:
		return fmt.Sprintf("builtin::%v", e.primitive)
	case ExpList:
		elems := make([]string, 0)
		for _, e := range e.list {
			elems = append(elems, e.String())
		}
		return fmt.Sprintf("(%v)", strings.Join(elems, " "))
	case ExpString:
		return e.string
	case ExpInteger:
		return fmt.Sprintf("%d", e.integer)
	case ExpFloat:
		return fmt.Sprintf("%v", e.float)
	case ExpSymbol:
		return e.symbol
	case ExpBool:
		return fmt.Sprintf("%v", e.bool)
	case ExpQuote:
		return e.quote.String()
	case ExpNil:
		return "nil"
	case ExpFunction:
		return fmt.Sprintf("fn<%v %v>", e.functionName, e.functionParams)
	default:
		return fmt.Sprintf("unknown→%#v", e)
	}
}

// DebugString provides type information for expressions
func (e Expression) DebugString() string {
	switch e.tag {
	case ExpPrimitive:
		return fmt.Sprintf("builtin::%v", e.primitive)
	case ExpList:
		elems := make([]string, 0)
		for _, e := range e.list {
			elems = append(elems, e.DebugString())
		}
		return fmt.Sprintf("(%v)", strings.Join(elems, " "))
	case ExpString:
		return "str→" + string(e.string)
	case ExpInteger:
		return fmt.Sprintf("int→%d", e.integer)
	case ExpFloat:
		return fmt.Sprintf("float→%f", e.float)
	case ExpSymbol:
		return "sym→" + string(e.symbol)
	case ExpBool:
		return fmt.Sprintf("bool→%v", e.bool)
	case ExpQuote:
		return "(quote " + e.quote.DebugString() + ")"
	case ExpNil:
		return "nil"
	default:
		return fmt.Sprintf("unknown→%#v", e)
	}
}

// IsSymbol returns true if expression is a symbol
func (e Expression) IsSymbol() bool {
	return e.tag == ExpSymbol
}

// IsAtom returns true if expression is not a list
func (e Expression) IsAtom() bool {
	return !e.IsList()
}

// IsList returns true if expression is a list
func (e Expression) IsList() bool {
	return e.tag == ExpList
}

// Size returns the number of elements in the expression, or 1 if it's
// not a list.
func (e Expression) Size() int {
	if e.IsList() {
		return len(e.list)
	}
	return 1
}

// IsPrimitive returns true if expression is builtin function.
func (e Expression) IsPrimitive() bool {
	return e.tag == ExpPrimitive
}

// IsFunction returns true of the expression represents a function
func (e Expression) IsFunction() bool {
	return e.tag == ExpFunction
}

// IsQuote returns true if expr is a quote
func (e Expression) IsQuote() bool {
	return e.tag == ExpQuote
}

// IsTruthy returns true of expr isn't false or nil.
func (e Expression) IsTruthy() bool {
	if e.tag == ExpNil {
		return false
	}

	if e.tag == ExpBool {
		return e.bool
	}

	return true
}

// InvokePrimitive returns the results of a functino application.
func (e Expression) InvokePrimitive(params []Expression) (Expression, error) {
	return e.primitive(params)
}

// Equals returns true of the values of e1 and e2 match
func (e Expression) Equals(e2 Expression) bool {
	return e.hash == e2.hash
}

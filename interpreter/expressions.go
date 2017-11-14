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
	typ       ExpressionType
	hash      uint32
	string    string
	integer   int64
	float     float64
	symbol    string
	bool      bool
	list      []Expression
	quote     *Expression
	primitive primitiveFunc
}

// NewExpr constructs a new expression of the given type
func NewExpr(typ ExpressionType, value interface{}) Expression {

	hash := fnv.New32()
	if typ == ExpList {
		for _, e := range value.([]Expression) {
			fmt.Fprint(hash, e.hash)
		}
	} else {
		fmt.Fprint(hash, value)
	}

	e := Expression{typ: typ, hash: hash.Sum32()}
	switch typ {
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
		log.Fatalf("unable to create new expr of type %v", typ)
	}
	return e
}

// NilExpression represents nil
var NilExpression = Expression{typ: ExpNil}

// TrueExpression for a boolean true
var TrueExpression = NewExpr(ExpBool, true)

// FalseExpression for a boolean false
var FalseExpression = NewExpr(ExpBool, false)

// StartsWith returns true if first elem in list is named prefix.
func (e Expression) StartsWith(prefix string) bool {
	if e.typ != ExpList {
		return false
	}

	elem := e.Head()
	return prefix == elem.symbol
}

// Head returns the first element of the list.
func (e Expression) Head() Expression {
	if e.typ != ExpList {
		log.Fatalf("Can't take the head of a %v", e)
	}

	if len(e.list) == 0 {
		return NilExpression
	}
	return e.list[0]
}

// Tail returns the rest of the elements of a list.
func (e Expression) Tail() Expression {
	if e.typ != ExpList {
		log.Fatalf("Can't take the tail of a %v", e)
	}

	if len(e.list) == 0 {
		return e
	}
	return Expression{typ: ExpList, list: e.list[1:]}
}

func (e Expression) String() string {
	switch e.typ {
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
	default:
		return fmt.Sprintf("unknown→%#v", e)
	}
}

// DebugString provides type information for expressions
func (e Expression) DebugString() string {
	switch e.typ {
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
	return e.typ == ExpSymbol
}

// IsAtom returns true if expression is not a list
func (e Expression) IsAtom() bool {
	return !e.IsList()
}

// IsList returns true if expression is a list
func (e Expression) IsList() bool {
	return e.typ == ExpList
}

func (e Expression) Size() int {
	if e.IsList() {
		return len(e.list)
	}
	return 1
}

// IsPrimitive returns true if expression is builtin function.
func (e Expression) IsPrimitive() bool {
	return e.typ == ExpPrimitive
}

// IsQuote returns true if expr is a quote
func (e Expression) IsQuote() bool {
	return e.typ == ExpQuote
}

// IsTruthy returns true of expr isn't false or nil.
func (e Expression) IsTruthy() bool {
	if e.typ == ExpNil {
		return false
	}

	if e.typ == ExpBool {
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

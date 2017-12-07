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
	"hash/fnv"
	"log"
	"os"
	"strings"
)

// ExpressionType is the type of expression
type ExpressionType int

// ExpressionTypes
const (
	ExpNil ExpressionType = iota
	ExpPrimitive
	ExpFunction
	ExpLambda
	ExpList
	ExpString
	ExpInteger
	ExpFloat
	ExpSymbol
	ExpBool
	ExpQuote
	ExpFile    // represents a file-handle
	ExpHashMap // 12
	ExpThunk   // 13
)

// ExprTypeName returns the type name of an expression type
func ExprTypeName(v ExpressionType) string {
	names := map[ExpressionType]string{
		ExpNil:       "nil",
		ExpPrimitive: "primitive",
		ExpFunction:  "function",
		ExpLambda:    "lambda",
		ExpList:      "list",
		ExpString:    "string",
		ExpInteger:   "integer",
		ExpFloat:     "float",
		ExpSymbol:    "symbol",
		ExpBool:      "bool",
		ExpQuote:     "quote",
		ExpFile:      "file",
		ExpHashMap:   "hash-map",
		ExpThunk:     "thunk",
	}

	value, ok := names[v]
	if !ok {
		return "unknown"
	}
	return value
}

// Type returns the type name of an expression
func (e Expression) Type() string {
	return ExprTypeName(e.tag)
}

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
	functionEnv    *Environment
	file           *fileData
	hashMap        *HakiHashMap
	thunkValue     *Expression
}

func hashIt(values ...interface{}) uint32 {
	hash := fnv.New32()
	for _, v := range values {
		fmt.Fprint(hash, v)
	}
	return hash.Sum32()
}

var genSymCounter = 0

// GenSym produces a unique symbol name per runtime.
func GenSym(prefix string) Expression {
	genSymCounter = genSymCounter + 1
	return NewExpr(ExpSymbol, fmt.Sprintf("%v%v", prefix, genSymCounter))
}

// WrapImplicitDo wraps an expression in a do block
func WrapImplicitDo(body []Expression) Expression {
	do := NewExpr(ExpSymbol, "do")
	return NewExpr(ExpList, append([]Expression{do}, body...))
}

// WithLambda returns an expression wrapped in a lambda
func WithLambda(env *Environment, body Expression) Expression {
	fn := NewExpr(ExpSymbol, "fn")
	params := NewExpr(ExpList, []Expression{})
	result := NewLambdaExpr(env, fn, params, body)
	return result
}

// NewStringExpr returns an expression representing a string
func NewStringExpr(s string) Expression {
	return NewExpr(ExpString, s)
}

// NewBoolExpr returns an expression representing a boolean value
func NewBoolExpr(b bool) Expression {
	return NewExpr(ExpBool, b)
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

// NewLambdaExpr produces a new lambda expression (folding over environment)
func NewLambdaExpr(env *Environment, name, params, body Expression) Expression {
	p := params
	b := body
	e := env.Clone()
	return Expression{
		tag:            ExpLambda,
		hash:           hashIt(ExpLambda, name, p.hash, b.hash, e),
		functionName:   name.symbol,
		functionParams: &p,
		functionBody:   &b,
		functionEnv:    e,
	}
}

// NewThunkExpr returns a thunk expression, a non-closed over lambda for let bindings.
func NewThunkExpr(body Expression) Expression {
	n := GenSym("t")
	b := body
	return Expression{
		tag:          ExpThunk,
		hash:         hashIt(ExpThunk, n.hash, b.hash),
		functionName: n.symbol,
		functionBody: &b,
		thunkValue:   nil,
	}
}

// NewExpr constructs a new expression of the given type
func NewExpr(tag ExpressionType, value interface{}) Expression {

	// The complicated types are handled elsewhere.
	switch tag {

	case ExpList:
		return NewListExpr(value.([]Expression))

	case ExpHashMap:
		return NewHashMapExpr(value.(*HakiHashMap))

	case ExpFile:
		return NewFileHandleExpr(value.(*os.File))
	}

	// Simple types here.

	e := Expression{tag: tag, hash: hashIt(tag, value)}
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
		log.Fatalf("unable to create new expr of type %v", ExprTypeName(tag))
	}
	return e
}

// NilExpression represents nil
var NilExpression = Expression{tag: ExpNil}

// TrueExpression for a boolean true
var TrueExpression = NewBoolExpr(true)

// FalseExpression for a boolean false
var FalseExpression = NewBoolExpr(false)

// StdinExpression represents the STDIN file handle
var StdinExpression = NewFileHandleExpr(os.Stdin)

// StdoutExpression represents the STDOUT file handle
var StdoutExpression = NewFileHandleExpr(os.Stdout)

// StderrExpression represents the STDERR file handle
var StderrExpression = NewFileHandleExpr(os.Stderr)

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
		return "\"" + e.string + "\""
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
	case ExpLambda:
		return fmt.Sprintf("lambda<%v|%v %v>", e.functionName, e.functionParams, e.functionBody)
	case ExpFunction:
		return fmt.Sprintf("fn<%v %v>", e.functionName, e.functionParams)
	case ExpThunk:
		return fmt.Sprintf("thunk<%v %v>", e.functionName, e.functionBody)
	case ExpFile:
		status := " (closed)"
		if e.file.isOpen {
			status = ""
		}
		return "file://" + e.file.path + status
	case ExpHashMap:
		return e.hashMap.String()
	default:
		return fmt.Sprintf("unknown→%#v", e)
	}
}

// IsNil returns true if the expression represents nil
func (e Expression) IsNil() bool {
	return e.tag == ExpNil
}

// IsBool returns true if the expression represents nil
func (e Expression) IsBool() bool {
	return e.tag == ExpBool
}

// IsSymbol returns true if expression is a symbol
func (e Expression) IsSymbol() bool {
	return e.tag == ExpSymbol
}

// IsInteger returns true if expression is an integer
func (e Expression) IsInteger() bool {
	return e.tag == ExpInteger
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

// IsFunction returns true if the expression represents a function
func (e Expression) IsFunction() bool {
	return e.tag == ExpFunction
}

// IsLambda returns true if the expression represents an anonymous
// function value
func (e Expression) IsLambda() bool {
	return e.tag == ExpLambda
}

// IsThunk is true of expression is a thunk
func (e Expression) IsThunk() bool {
	return e.tag == ExpThunk
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

// IsInvokable if it can be called
func (e Expression) IsInvokable() bool {
	return e.tag == ExpFunction || e.tag == ExpPrimitive || e.tag == ExpLambda
}

// InvokePrimitive returns the results of a functino application.
func (e Expression) InvokePrimitive(params []Expression) (Expression, error) {
	return e.primitive(params)
}

// Equals returns true of the values of e1 and e2 match
func (e Expression) Equals(e2 Expression) bool {
	return e.hash == e2.hash
}

// IntValue of the expression
func (e Expression) IntValue() int64 {
	return e.integer
}

// Testing help

// Value returns the raw golang value of the expression
func (e Expression) Value() interface{} {
	switch e.tag {
	case ExpPrimitive:
		return fmt.Sprintf("builtin::%v", e.primitive)
	case ExpList:
		elems := make([]interface{}, 0)
		for _, e := range e.list {
			elems = append(elems, e.Value())
		}
		return elems
	case ExpString:
		return e.string
	case ExpInteger:
		return e.integer
	case ExpFloat:
		return e.float
	case ExpSymbol:
		return e.symbol
	case ExpBool:
		return e.bool
	case ExpQuote:
		return e.quote.String()
	case ExpNil:
		return nil
	case ExpLambda:
		return fmt.Sprintf("lambda<%v|%v %v>", e.functionName, e.functionParams, e.functionBody)
	case ExpFunction:
		return fmt.Sprintf("fn<%v %v>", e.functionName, e.functionParams)
	case ExpHashMap:
		return e.hashMap
	default:
		return fmt.Sprintf("unknown→%#v", e)
	}
}

// IsEqual to expression
func (e Expression) IsEqual(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v == e.bool
	case int:
		return int(e.integer) == v
	case int64:
		return e.integer == v
	case float64:
		return e.float == v
	case string:
		return e.string == v || e.symbol == v
	case []string:
		exprs := make([]Expression, 0)
		for _, s := range v {
			exprs = append(exprs, NewExpr(ExpString, s))
		}
		return e.Equals(NewExpr(ExpList, exprs))
	case []int64:
		// TODO more efficent to convert e to ints, than this
		exprs := make([]Expression, 0)
		for _, i := range v {
			exprs = append(exprs, NewExpr(ExpInteger, i))
		}
		return e.Equals(NewExpr(ExpList, exprs))
	default:
		fmt.Printf("equality value '%v' of unknown type\n", value)
		return false
	}
}

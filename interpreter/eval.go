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
	"errors"
	"fmt"
)

func nilExpr(reason string, params ...interface{}) (Expression, error) {
	return NilExpression, fmt.Errorf(reason, params...)
}

func apply(env *Environment, op Expression, args []Expression) (Expression, error) {

	theOp, err := Evaluate(env, op)
	if err != nil {
		return NilExpression, err
	}

	argv := make([]Expression, 0)
	for _, a := range args {
		param, err := Evaluate(env, a)
		if err != nil {
			return NilExpression, err
		}
		argv = append(argv, param)
	}

	// Primitive operator
	if theOp.IsPrimitive() {
		return theOp.InvokePrimitive(argv)
	}

	argc := len(args)
	paramc := len(theOp.functionParams.list)
	if argc != paramc {
		return nilExpr("function '%v' takes %v param(s), you provided %v", theOp.functionName, paramc, argc)
	}

	// Global function
	if theOp.IsFunction() {
		newEnv := env.ExtendEnvironment(*theOp.functionParams, argv)
		return Evaluate(newEnv, *theOp.functionBody)
	}

	// Anonymous (lambda) function
	if theOp.IsLambda() {
		newEnv := theOp.functionEnv.ExtendEnvironment(*theOp.functionParams, argv)
		return Evaluate(newEnv, *theOp.functionBody)
	}

	return nilExpr("function not found: '%v'", theOp)
}

func evalIf(env *Environment, exprs Expression) (Expression, error) {
	argc := len(exprs.list)
	if argc < 2 {
		return NilExpression, fmt.Errorf("too few arguments (%v) to if", argc)
	}
	if argc > 3 {
		return NilExpression, fmt.Errorf("too many arguments (%v) to if", argc)
	}

	test, err := Evaluate(env, exprs.list[0])
	if err != nil {
		return NilExpression, err
	}

	var result Expression

	if test.IsTruthy() {
		result, err = Evaluate(env, exprs.list[1])
		if err != nil {
			return NilExpression, err
		}
	} else if argc >= 3 {
		result, err = Evaluate(env, exprs.list[2])
		if err != nil {
			return NilExpression, err
		}
	}
	return result, nil
}

func evalDo(env *Environment, exprs Expression) (Expression, error) {

	if !exprs.IsList() {
		return Evaluate(env, exprs)
	}

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}
	}
	return result, err
}

func evalLet(env *Environment, clauses Expression, body Expression) (Expression, error) {

	if !clauses.IsList() {
		return nilExpr("let bindings should be a list (let (a 1 b 2) ...)")
	}

	if !(clauses.Size()%2 == 0) {
		return nilExpr("let bindings must contain an even number of left/right pairs")
	}

	params := make([]Expression, 0)
	args := make([]Expression, 0)

	for i := 0; i < clauses.Size(); i = i + 2 {
		param := clauses.list[i]
		arg, err := Evaluate(env, clauses.list[i+1])
		if err != nil {
			return NilExpression, err
		}

		params = append(params, param)
		args = append(args, arg)
	}

	newEnv := env.ExtendEnvironment(NewExpr(ExpList, params), args)
	doBlock := WrapImplicitDo(body.list)

	return Evaluate(newEnv, doBlock)
}

func evalAnd(env *Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if !result.IsTruthy() {
			return result, nil
		}
	}

	return result, nil
}

func evalOr(env *Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if result.IsTruthy() {
			return result, nil
		}
	}
	return result, nil
}

func evalDef(env *Environment, name Expression, body Expression) (Expression, error) {

	if !name.IsSymbol() {
		return NilExpression, errors.New("def name must be a symbol")
	}

	do := WrapImplicitDo(body.list)
	value, err := Evaluate(env, do)

	if err != nil {
		return NilExpression, err
	}

	env.Set(name, value)
	return value, nil
}

func evalFunction(env *Environment, name Expression, params Expression, body Expression) (Expression, error) {

	if !name.IsSymbol() {
		return nilExpr("defun name ← name must be a symbol")
	}

	if !params.IsList() {
		return nilExpr("defun name (params) ← parameters must be a list")
	}

	f := NewFunctionExpr(name, params, body)
	env.Set(name, f)
	return f, nil
}

func evalLambda(env *Environment, params Expression, body Expression) (Expression, error) {

	if !params.IsList() {
		return nilExpr("in (fn (params) (body)) ← params must be a list")
	}

	name := GenSym("fn") // necessary?
	f := NewLambdaExpr(env, name, params, body.Head())
	return f, nil
}

// Evaluate an expression
func Evaluate(env *Environment, expr Expression) (Expression, error) {

	if expr.IsSymbol() {
		found, value := env.Lookup(expr.symbol)
		if !found {
			return nilExpr("value not found for '%v'", expr.String())
		}
		return value, nil
	}

	if expr.IsQuote() {
		return *expr.quote, nil
	}

	if expr.IsAtom() {
		return expr, nil
	}

	if expr.IsList() {
		if expr.StartsWith("do") {
			return evalDo(env, expr.Tail())
		}
		if expr.StartsWith("let") {
			return evalLet(env, expr.Tail().Head(), expr.Tail().Tail())
		}
		if expr.StartsWith("and") {
			return evalAnd(env, expr.Tail())
		}
		if expr.StartsWith("or") {
			return evalOr(env, expr.Tail())
		}
		if expr.StartsWith("if") {
			return evalIf(env, expr.Tail())
		}
		if expr.StartsWith("def") {
			def := expr.Tail()
			return evalDef(env, def.Head(), def.Tail())
		}
		if expr.StartsWith("defun") {
			name := expr.Tail().Head()
			params := expr.Tail().Tail().Head()
			body := expr.Tail().Tail().Tail()
			return evalFunction(env, name, params, body)
		}
		if expr.StartsWith("fn") {
			params := expr.Tail().Head()
			body := expr.Tail().Tail()
			return evalLambda(env, params, body)
		}
		return apply(env, expr.Head(), expr.Tail().list)
	}

	return NilExpression, fmt.Errorf("unable to eval expression [%v]", expr)
}

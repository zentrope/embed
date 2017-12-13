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
)

func (x TcoInterpreter) evalList(env *Environment, args []Expression) ([]Expression, error) {
	argv := make([]Expression, 0)
	for _, a := range args {
		param, err := x.Evaluate(env, a)
		if err != nil {
			return []Expression{}, err
		}
		argv = append(argv, param)
	}
	return argv, nil
}

func isValidArity(fn Expression, args []Expression) (bool, error) {

	if !fn.IsInvokable() {
		return false, fmt.Errorf("fn '%v' (%v) is not invokable", fn, fn.Type())
	}

	argc := len(args)
	paramc := len(fn.functionParams.list)
	if argc == paramc {
		return true, nil
	}

	return false, fmt.Errorf("fn '%v' takes %v param(s), you provided %v",
		fn.functionName, paramc, argc)
}

//-----------------------------------------------------------------------------
// LOOP
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalLoop(env *Environment, args Expression) (Expression, error) {

	if err := typeCheck("(loop (fn (x) ...) lst)", args.list,
		ckArity(2), ckFuncable(0)); err != nil {
		return NilExpression, err
	}

	fn := args.list[0]

	lst, err := x.Evaluate(env, args.list[1])
	if err != nil {
		return NilExpression, err
	}

	for _, e := range lst.list {
		_, err := x.Evaluate(env, NewListExpr([]Expression{fn, e}))
		if err != nil {
			return NilExpression, err
		}
	}
	return NilExpression, nil
}

//-----------------------------------------------------------------------------
// LOOP-INDEX
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalLoopIndex(env *Environment, args Expression) (Expression, error) {

	if err := typeCheck("(loop-index (fn (i x) ...) lst)", args.list,
		ckArity(2), ckFuncable(0)); err != nil {
		return NilExpression, err
	}

	fn := args.list[0]

	lst, err := x.Evaluate(env, args.list[1])
	if err != nil {
		return NilExpression, err
	}

	for i, e := range lst.list {
		_, err := x.Evaluate(env, NewListExpr([]Expression{fn, NewIntExpr(int64(i)), e}))
		if err != nil {
			return NilExpression, err
		}
	}
	return NilExpression, nil
}

//-----------------------------------------------------------------------------
// IF
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalIf(env *Environment, exprs Expression) (Expression, error) {
	argc := len(exprs.list)
	if argc < 2 {
		return NilExpression, fmt.Errorf("too few arguments (%v) to if", argc)
	}
	if argc > 3 {
		return NilExpression, fmt.Errorf("too many arguments (%v) to if", argc)
	}

	test, err := x.Evaluate(env, exprs.list[0])
	if err != nil {
		return NilExpression, err
	}

	if test.IsTruthy() {
		return exprs.list[1], nil
	}

	if argc >= 3 {
		return exprs.list[2], nil
	}

	return NilExpression, nil
}

//-----------------------------------------------------------------------------
// LET
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalLet(env *Environment, clauses, body Expression) (*Environment, Expression, error) {

	if !clauses.IsList() {
		return env, NilExpression,
			errors.New("let bindings should be a list (let (a 1 b 2) ...)")
	}

	if !(clauses.Size()%2 == 0) {
		return env, NilExpression,
			errors.New("let bindings must contain an even number of left/right pairs")
	}

	if clauses.Size() == 0 {
		return env, body, nil
	}

	bindingNames := make([]Expression, 0)
	bindingVals := make([]Expression, 0)

	for i := 0; i < clauses.Size(); i += 2 {
		name := clauses.list[i]
		val := clauses.list[i+1]
		bindingNames = append(bindingNames, name)
		bindingVals = append(bindingVals, NewThunkExpr(val))
	}

	newEnv := env.ExtendEnvironment(NewListExpr(bindingNames), bindingVals)

	doBlock := WrapImplicitDo(body.list)
	return newEnv, doBlock, nil
}

//-----------------------------------------------------------------------------
// AND
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalAnd(env *Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = x.Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if !result.IsTruthy() {
			return result, nil
		}
	}

	return result, nil
}

//-----------------------------------------------------------------------------
// OR
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalOr(env *Environment, exprs Expression) (Expression, error) {

	var result Expression
	var err error

	for _, e := range exprs.list {
		result, err = x.Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}

		if result.IsTruthy() {
			return result, nil
		}
	}
	return result, nil
}

//-----------------------------------------------------------------------------
// DO
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalDo(env *Environment, exprs Expression) (Expression, error) {

	if !exprs.IsList() {
		return exprs, nil
	}

	for _, e := range exprs.list[:len(exprs.list)-1] {
		_, err := x.Evaluate(env, e)
		if err != nil {
			return NilExpression, err
		}
	}
	return exprs.list[len(exprs.list)-1], nil
}

//-----------------------------------------------------------------------------
// DEF
//-----------------------------------------------------------------------------

func (x TcoInterpreter) evalDef(env *Environment, name, body Expression) (Expression, error) {

	if !name.IsSymbol() {
		return nilExpr("(def name expr) «- name must be a symbol, not '%v'", ExprTypeName(name.tag))
	}

	argc := len(body.list)

	if argc == 0 {
		env.Set(name, NIL)
		return NIL, nil
	}

	if argc != 1 {
		return nilExpr("(def name expr) «-- takes 1 expr, you offered %v", argc)
	}

	do := WrapImplicitDo(body.list)
	value, err := x.Evaluate(env, do)

	if err != nil {
		return NilExpression, err
	}

	env.Set(name, value)
	return value, nil
}

func (x TcoInterpreter) evalDefun(env *Environment, name, params Expression, body Expression) (Expression, error) {

	sig := "(defun name (params) body…)"

	if !name.IsSymbol() {
		return nilExpr("%v «-- name must be a symbol", sig)
	}

	if !params.IsList() {
		return nilExpr("%v «-- params must be a list", sig)
	}

	f := NewFunctionExpr(env, name, params, body)
	env.Set(name, f)
	return f, nil
}

func (x TcoInterpreter) evalLambda(env *Environment, params, body Expression) (Expression, error) {
	if !params.IsList() {
		return nilExpr("in (fn (params) (body)) ← params must be a list")
	}

	f := NewLambdaExpr(env, GenSym("fn"), params, body.Head())
	return f, nil
}

//-----------------------------------------------------------------------------
// EVAL
//-----------------------------------------------------------------------------

// Evaluate an expression in an environment, returning an expression.
func (x TcoInterpreter) Evaluate(env *Environment, expr Expression) (Expression, error) {
	var err error

	for {
		switch expr.tag {

		case ExpNil:
			return expr, nil

		case ExpSymbol:
			found, value := env.Lookup(expr.symbol)
			if !found {
				return nilExpr("value not found for '%v'", expr.String())
			}

			if value.IsThunk() {
				boundValue, err := x.Evaluate(env, *value.functionBody)
				if err != nil {
					return NilExpression, err
				}
				env.Replace(expr, boundValue)
				return boundValue, nil
			}
			return value, nil

		case ExpQuote:
			return *expr.quote, nil

		case ExpInteger, ExpString, ExpFloat, ExpBool:
			return expr, nil

		case ExpList:
			first := expr.Head()
			rest := expr.Tail()

			switch first.symbol {

			case "loop":
				return x.evalLoop(env, rest)

			case "loop-index":
				return x.evalLoopIndex(env, rest)

			case "if":
				expr, err = x.evalIf(env, rest)
				if err != nil {
					return expr, err
				}

			case "and":
				return x.evalAnd(env, rest)

			case "or":
				return x.evalOr(env, rest)

			case "do":
				expr, err = x.evalDo(env, rest)
				if err != nil {
					return NilExpression, err
				}

			case "let":
				env, expr, err = x.evalLet(env, rest.Head(), rest.Tail())
				if err != nil {
					return NilExpression, err
				}

			case "def":
				return x.evalDef(env, rest.Head(), rest.Tail())

			case "defun":
				name := rest.Head()
				params := rest.Tail().Head()
				body := rest.Tail().Tail()
				return x.evalDefun(env, name, params, body)

			case "fn", "lambda":
				params := rest.Head()
				body := rest.Tail()
				return x.evalLambda(env, params, body)

			default: // apply
				fn, err := x.Evaluate(env, first)
				if err != nil {
					return NilExpression, err
				}

				if !fn.IsInvokable() {
					println("Not invokable.")
				}
				argv, err := x.evalList(env, rest.list)
				if err != nil {
					return NilExpression, err
				}

				if fn.IsPrimitive() {
					ret, err := fn.InvokePrimitive(argv)
					return ret, err
				}

				ok, err := isValidArity(fn, rest.list)
				if !ok {
					return NilExpression, err
				}

				if fn.IsLambda() {
					env = fn.functionEnv.ExtendEnvironment(*fn.functionParams, argv)
					expr = *fn.functionBody
				} else if fn.IsFunction() {
					env = fn.functionEnv.ExtendEnvironment(*fn.functionParams, argv)
					// env = env.ExtendEnvironment(*fn.functionParams, argv)
					expr = *fn.functionBody
				} else {
					return nilExpr("unable to apply %v", fn)
				}
			}
		}

	} // for

}

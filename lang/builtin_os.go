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
	"os"
	"strings"
)

var osBuiltins = primitivesMap{
	"env":         _env,
	"environment": _environment,
}

func _env(args []Expression) (Expression, error) {
	if err := typeCheck("(env string)", args,
		ckArityAtLeast(1), ckString(0), ckOptString(1)); err != nil {
		return NilExpression, err
	}

	name := args[0].string

	result := os.Getenv(name)
	if result == "" {
		if len(args) == 2 {
			return NewStringExpr(args[1].string), nil
		}
		return NilExpression, nil
	}
	return NewStringExpr(result), nil
}

func _environment(args []Expression) (Expression, error) {

	if err := typeCheck("(environment)", args, ckArity(0)); err != nil {
		return NilExpression, err
	}

	env := os.Environ()
	results := newHakiMap()

	for _, e := range env {
		words := strings.Split(e, "=")
		results.set(NewStringExpr(words[0]), NewStringExpr(words[1]))
	}

	return NewHashMapExpr(results), nil
}

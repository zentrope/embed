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

var logicBuiltins = primitivesMap{
	"=":   _equals,
	"not": _not,
}

func _not(args []Expression) (Expression, error) {
	if err := typeCheck("(not v)", args, ckArity(1)); err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpBool, !args[0].IsTruthy()), nil
}

func _equals(args []Expression) (Expression, error) {

	if err := typeCheck("(= v1 ... vn)", args, ckArityAtLeast(1)); err != nil {
		return NilExpression, err
	}

	sentinel := args[0]

	for _, a := range args[1:] {
		if !a.Equals(sentinel) {
			return FalseExpression, nil
		}
	}
	return TrueExpression, nil
}

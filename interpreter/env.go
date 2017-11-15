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

// Environment represents bindings
type Environment struct {
	data map[string]Expression
}

// NewEnvironment contains bindings
func NewEnvironment() *Environment {
	data := make(map[string]Expression, 0)

	for name, fn := range builtins {
		data[name] = NewExpr(ExpPrimitive, fn)
	}

	return &Environment{data: data}
}

// Lookup a value in the environment
func (env *Environment) Lookup(key string) (bool, Expression) {
	value := env.data[key]
	if value.tag == 0 {
		return false, NilExpression
	}
	return true, value
}

// Set a value in the current environment frame
func (env *Environment) Set(key Expression, value Expression) {
	env.data[key.symbol] = value
}

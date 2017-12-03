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

type primitiveFunc func(args []Expression) (Expression, error)
type primitivesMap map[string]primitiveFunc

var builtins = make(primitivesMap, 0)

func init() {
	prims := []primitivesMap{
		logicBuiltins,   // builtins_logic
		mathBuiltins,    // builtins_math
		stringBuiltins,  // builtins_string
		listBuiltins,    // builtins_list
		fileioBuiltins,  // builtins_fileio
		hashmapBuiltins, // builtins_hashmap
		writeBuiltins,   // builtins_write
	}
	for _, prim := range prims {
		for name, fn := range prim {
			builtins[name] = fn
		}
	}
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

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

// TODO: Move to type checking.
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

// TODO: Move to type checking.
func verifyNums(args []Expression) error {
	for _, arg := range args {
		switch arg.tag {
		case ExpInteger, ExpFloat:
			continue
		default:
			return errors.New("all arguments must be numbers")
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
// Type checking
//-----------------------------------------------------------------------------

type spec func(string, []Expression) error

func ckArity(numArgs int) spec {
	return func(sig string, args []Expression) error {
		if len(args) != numArgs {
			return fmt.Errorf("'%v' expects '%v' args, you provided '%v'",
				sig, numArgs, len(args))
		}
		return nil
	}
}

func ckArityAtLeast(numArgs int) spec {
	return func(sig string, args []Expression) error {
		if len(args) <= numArgs {
			return fmt.Errorf("'%v' expects at least '%v' args, you provided '%v'",
				sig, numArgs, len(args))
		}
		return nil
	}
}

func ckType(pos int, tag ExpressionType) spec {
	return func(sig string, args []Expression) error {
		if args[pos].tag != tag {
			return fmt.Errorf("'%v' expects arg %v to be type '%v', not '%v'",
				sig, pos+1, ExprTypeName(tag), ExprTypeName(args[pos].tag))
		}
		return nil
	}
}

func ckString(pos int) spec {
	return ckType(pos, ExpString)
}

func ckMap(pos int) spec {
	return ckType(pos, ExpHashMap)
}

func ckList(pos int) spec {
	return ckType(pos, ExpList)
}

func typeCheck(sig string, args []Expression, specs ...spec) error {
	for _, spec := range specs {
		if err := spec(sig, args); err != nil {
			return err
		}
	}
	return nil
}

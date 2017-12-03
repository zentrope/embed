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
	"strings"
)

var hashmapBuiltins = primitivesMap{
	"hget":    _hget,
	"hmap":    _hmap,
	"hset":    _hset,
	"hmap?":   _hmapP,
	"hkeys":   _hkeys,
	"hvals":   _hvals,
	"hget-in": _hgetin,
}

// HakiHashMap represents a hash-map type in the Haki language.
type HakiHashMap struct {
	keys map[uint32]Expression
	vals map[uint32]Expression
}

func newHakiMap() *HakiHashMap {
	return &HakiHashMap{
		keys: make(map[uint32]Expression),
		vals: make(map[uint32]Expression),
	}
}

func (hmap *HakiHashMap) set(key, value Expression) {
	hmap.keys[key.hash] = key
	hmap.vals[key.hash] = value
}

func (hmap *HakiHashMap) copy() *HakiHashMap {
	newMap := newHakiMap()
	for lookup := range hmap.keys {
		newMap.set(hmap.keys[lookup], hmap.vals[lookup])
	}
	return newMap
}

func (hmap *HakiHashMap) String() string {
	sections := make([]string, 0)
	for hash, key := range hmap.keys {
		sections = append(sections, fmt.Sprintf("%v: %v", key, hmap.vals[hash]))
	}

	return "(hmap " + strings.Join(sections, ", ") + ")"
}

// NewHashMapExpr returns an expression wrapper around a hash map
func NewHashMapExpr(hmap *HakiHashMap) Expression {
	data := make([]interface{}, 0)
	data = append(data, ExpHashMap)

	for _, val := range hmap.vals {
		data = append(data, val.hash)
	}

	return Expression{
		tag:     ExpHashMap,
		hash:    hashIt(data...),
		hashMap: hmap,
	}
}

//-----------------------------------------------------------------------------
// implementations
//-----------------------------------------------------------------------------

func _hmap(args []Expression) (Expression, error) {
	argc := len(args)
	sig := "(hmap k v ...)"

	if (argc % 2) != 0 {
		return nilExpr("%v expects an even number of params.", sig)
	}

	hmap := newHakiMap()
	for i := 0; i < argc; i += 2 {
		key := args[i]
		value := args[i+1]
		hmap.set(key, value)
	}

	return NewHashMapExpr(hmap), nil
}

func _hmapP(args []Expression) (Expression, error) {
	sig := "(hmap? val)"
	argc := len(args)
	if argc > 1 {
		return nilExpr("%v takes 1 arg, you provided %v", sig, argc)
	}

	return NewExpr(ExpBool, args[0].tag == ExpHashMap), nil
}

func _hkeys(args []Expression) (Expression, error) {
	sig := "(hkeys m)"
	argc := len(args)
	if argc != 1 {
		return nilExpr("%v expects at least 1 arg, you provided %v",
			sig, argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("%v expects 'm' to be a 'hash-map', not '%v'", sig,
			ExprTypeName(args[0].tag))
	}

	exprs := make([]Expression, 0)
	for _, v := range args[0].hashMap.keys {
		exprs = append(exprs, v)
	}

	return NewExpr(ExpList, exprs), nil
}

func _hvals(args []Expression) (Expression, error) {
	sig := "(hvals m)"
	argc := len(args)
	if argc != 1 {
		return nilExpr("%v expects at least 1 arg, you provided %v", sig, argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("%v expects 'm' to be a 'hash-map', not '%v'",
			sig,
			ExprTypeName(args[0].tag))
	}

	exprs := make([]Expression, 0)
	for _, v := range args[0].hashMap.vals {
		exprs = append(exprs, v)
	}

	return NewExpr(ExpList, exprs), nil
}

func _hget(args []Expression) (Expression, error) {
	sig := "(hget m k)"
	argc := len(args)
	if argc < 2 {
		return nilExpr("%v expects at least 2 args, you provided %v", sig, argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("%v expects 'm' to be a 'hash-map', not '%v'", sig,
			ExprTypeName(args[0].tag))
	}

	h := args[0].hashMap
	key := args[1].hash

	return h.vals[key], nil
}

func _hset(args []Expression) (Expression, error) {
	sig := "(hset m k v ...)"
	argc := len(args)
	if argc < 3 {
		return nilExpr("%v expects at least 3 args, you provided %v", sig,
			argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("%v expects 'm' to be a 'hash-map', not '%v'", sig,
			ExprTypeName(args[0].tag))
	}

	original := args[0]

	keyValues := args[1:]
	if (len(keyValues) % 2) != 0 {
		return nilExpr("(hset m k v ...) expects an even number of k/v params, you provided %v.",
			len(keyValues))
	}

	newMap := original.hashMap.copy()

	for i := 0; i < len(keyValues); i += 2 {
		key := keyValues[i]
		value := keyValues[i+1]
		newMap.set(key, value)
	}

	return NewHashMapExpr(newMap), nil
}

func _hgetin(args []Expression) (Expression, error) {
	sig := "(hget-in m [k & ks])"

	if err := runCheckers(sig, args, ckArity(2), ckMap(0), ckList(1)); err != nil {
		return NilExpression, err
	}

	m := args[0]

	for _, k := range args[1].list {
		if m.tag == ExpHashMap {
			m = m.hashMap.vals[k.hash]
		} else {
			return NilExpression, nil
		}
	}
	return m, nil
}

//-----------------------------------------------------------------------------
// Type checking (used in here as a test case)
//-----------------------------------------------------------------------------

type ckFn func(string, []Expression) error

func ckArity(numArgs int) ckFn {
	return func(sig string, args []Expression) error {
		if len(args) != numArgs {
			return fmt.Errorf("'%v' expects %v args, you provided %v", sig, numArgs, len(args))
		}
		return nil
	}
}

func ckType(pos int, tag ExpressionType) ckFn {
	return func(sig string, args []Expression) error {
		if args[pos].tag != tag {
			return fmt.Errorf("'%v' expects arg %v to be type '%v', not '%v'",
				sig, pos+1, ExprTypeName(tag), ExprTypeName(args[pos].tag))
		}
		return nil
	}
}

func ckMap(pos int) ckFn {
	return ckType(pos, ExpHashMap)
}

func ckList(pos int) ckFn {
	return ckType(pos, ExpList)
}

func runCheckers(sig string, args []Expression, checks ...ckFn) error {
	for _, check := range checks {
		if err := check(sig, args); err != nil {
			return err
		}
	}
	return nil
}

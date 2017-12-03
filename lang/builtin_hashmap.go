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

func _hmap(args []Expression) (Expression, error) {
	argc := len(args)
	if (argc % 2) != 0 {
		return nilExpr("(hmap k v ...) requires an even number of params.")
	}

	hmap := newHakiMap()
	for i := 0; i < argc; i += 2 {
		key := args[i]
		value := args[i+1]
		hmap.set(key, value)
	}

	return NewHashMapExpr(hmap), nil
}

func _hget(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 2 {
		return nilExpr("(hget m k) requires at least 2 args, you provided %v",
			argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("(hget m k) requires 'm' to be a 'hash-map', not '%v'",
			ExprTypeName(args[0].tag))
	}

	h := args[0].hashMap
	key := args[1].hash

	return h.vals[key], nil
}

func _hset(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 3 {
		return nilExpr("(hset m k v ...) requires at least 3 args, you provided %v",
			argc)
	}

	if args[0].tag != ExpHashMap {
		return nilExpr("(hset m k v ...) requires 'm' to be a 'hash-map', not '%v'",
			ExprTypeName(args[0].tag))
	}

	original := args[0]

	keyValues := args[1:]
	if (len(keyValues) % 2) != 0 {
		return nilExpr("(hset m k v ...) requires an even number of k/v params, you provided %v.",
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

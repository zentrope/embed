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
	"fmt"
	"strings"
)

// Sexp represents a parse s-expression
type Sexp interface {
	AsString() string
}

type sexpPrimitive primitiveFunc
type sexpList []Sexp
type sexpString string
type sexpInteger int64
type sexpFloat float64
type sexpSymbol string
type sexpBool bool
type sexpQuote struct {
	quote Sexp
}

// StartsWith returns true if first elem in list is named prefix.
func (s sexpList) StartsWith(prefix string) bool {
	elem := s.Head()
	switch t := elem.(type) {
	case sexpSymbol:
		return prefix == string(t)
	default:
		return false
	}
}

// Head returns the first term of the list.
func (s sexpList) Head() Sexp {
	if len(s) == 0 {
		return nil
	}
	return s[0]
}

// Tail returns rest of the list after the head is removed.
func (s sexpList) Tail() []Sexp {
	if len(s) == 0 {
		return s
	}
	return s[1:]
}

func (s sexpBool) AsString() string {
	return fmt.Sprintf("bool::%v", s)
}

func (s sexpPrimitive) AsString() string {
	return fmt.Sprintf("builtin::%v", s)
}

func (s sexpQuote) AsString() string {
	return "(quote " + s.quote.AsString() + ")"
}

func (s sexpList) AsString() string {
	elems := make([]string, 0)
	for _, e := range s {
		elems = append(elems, e.AsString())
	}
	return fmt.Sprintf("(%v)", strings.Join(elems, " "))
}

func (s sexpString) AsString() string {
	return "str::" + string(s)
}

func (s sexpInteger) AsString() string {
	return fmt.Sprintf("int::%d", s)
}

func (s sexpFloat) AsString() string {
	return fmt.Sprintf("float::%f", s)
}

func (s sexpSymbol) AsString() string {
	return "sym::" + string(s)
}

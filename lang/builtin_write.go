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

var writeBuiltins = primitivesMap{
	"prn": _prn,
}

func unquoteAll(value string) string {
	// Surely there's a better or builtin way to do this.
	reps := map[string]string{
		"\\n": "\n",
		"\\r": "\r",
		"\\t": "\t",
	}
	for k, v := range reps {
		value = strings.Replace(value, k, v, -1)
	}
	return value
}

func _prn(args []Expression) (Expression, error) {
	if len(args) == 0 {
		fmt.Println("")
		return NilExpression, nil
	}

	values := make([]string, 0)
	var value string
	for _, a := range args {
		value = a.String()
		if a.tag == ExpString {
			value = unquoteAll(a.string)
		}
		values = append(values, value)
	}
	fmt.Println(strings.Join(values, " "))
	return NilExpression, nil
}

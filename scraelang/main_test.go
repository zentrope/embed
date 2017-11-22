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

package scraelang

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/zentrope/embed/interpreter"
)

var NIL = interpreter.NilExpression

var flatre = regexp.MustCompile(`\s+`)

func flatten(s string) string {
	return strings.TrimSpace(flatre.ReplaceAllString(s, " "))
}

func evalForm(form string) (interpreter.Expression, error) {

	lang := interpreter.NewInterpreter(interpreter.TCO)
	env := interpreter.NewEnvironment()
	reader := interpreter.NewReader()
	reader.Append(interpreter.Core)

	doEval := func(form string) (interpreter.Expression, error) {
		tokens, err := interpreter.Tokenize(form)
		if err != nil {
			return NIL, err
		}

		parser := interpreter.NewParser(tokens)
		expr, err := parser.Parse()

		if err != nil {
			return NIL, err
		}

		return lang.Evaluate(env, expr)
	}

	for {
		if reader.IsBalanced() {
			form, err := reader.GetNextForm()
			if err != nil {
				if err == interpreter.ErrEOF {
					break
				}
				fmt.Printf(" ERROR: %v\n", err)
			}
			doEval(form)
			if form == "" {
				break
			}
		}
	}

	return doEval(form)
}

type form struct {
	tag      string
	expected interface{}
	form     string
}

func TestLetRecursive(t *testing.T) {
	table := []form{
		{"integer", int64(2), `(let (a 1 b (+ a 1)) b)`},
		{"integer", int64(0), `(let (a (fn (x) (if (= x 0) x (a (- x 1))))) (a 13))`},
		{"integer", int64(3), `(let (x 1) (let (a 2 b (fn () (+ a x))) (b))) `},
		{"list", []int64{1, 3, 5, 7, 9}, `(filter (fn (x) (odd? x)) (range 10))`},
	}
	for _, row := range table {
		t.Logf("letrec: %v", row.form)
		rc, err := evalForm(row.form)
		if err != nil {
			t.Error(err)
		}

		if !(rc.Type() == row.tag) {
			t.Errorf("Expected '%v' result: %v → %v → %v", row.tag, row.expected, row.form, rc)
		}

		if !rc.IsEqual(row.expected) {
			t.Errorf("Expected '%v' result: %v → %v → %v (%v)",
				row.tag, row.form, row.expected, rc, rc.Type())
		}
	}

}

func TestSimpleMath(t *testing.T) {

	table := []form{
		{"integer", int64(5), `(+ 2 3)`},
		{"integer", int64(1), `(- 100 99)`},
		{"integer", int64(-23), `(- 100 99 24)`},
		{"float", float64(2.1), `(+ 2 0.1)`},
		{"integer", int64(10), `(+ 1 (+ 2 6) (- 10 9))`},
	}

	for _, row := range table {
		t.Logf("math: %v", row.form)
		rc, err := evalForm(row.form)
		if err != nil {
			t.Error(err)
		}

		if !(rc.Type() == row.tag) {
			t.Errorf("Expected '%v' result: %v → %v → %v", row.tag, row.expected, row.form, rc)
		}

		if !rc.IsEqual(row.expected) {
			t.Errorf("Expected '%v' result: %v → %v → %v (%v)",
				row.tag, row.form, row.expected, rc, rc.Type())
		}
	}

}

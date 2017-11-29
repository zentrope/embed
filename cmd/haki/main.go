//
// Copyright (C) 2017 Keith Irwin
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

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	haki "github.com/zentrope/haki/lang"
)

func eval(interpreter haki.Interpreter, form string) haki.Expression {
	result, err := interpreter.Execute(form)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return haki.NilExpression
	}
	return result
}

const promptRepl = "repl> "
const promptMore = "   +> "

var flatre = regexp.MustCompile(`\s+`)

func flatten(s string) string {
	return strings.TrimSpace(flatre.ReplaceAllString(s, " "))
}

func readAll(interpreter haki.Interpreter, reader *haki.Reader, coreload bool) {
	for {
		if reader.IsBalanced() {
			form, err := reader.GetNextForm()
			if err != nil {
				if err == haki.ErrEOF {
					break
				}
				fmt.Printf(" ERROR: %v\n", err)
			}
			if form == "" {
				break
			}

			if coreload {
				fmt.Print(".")
				eval(interpreter, form)
			} else {
				fmt.Printf("%v\n", eval(interpreter, form))
			}
			continue
		}
	}
}

func main() {
	fmt.Println("Haki Repl")

	rl, err := readline.New(promptRepl)
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	// interpreter := haki.NewInterpreter(haki.Naive)
	// fmt.Println("Naive interpreter")
	interpreter := haki.NewInterpreter(haki.TCO)
	fmt.Println("* using tco interpreter")

	reader := haki.NewReader()

	// load core

	fmt.Print("* loading core")
	reader.Append(haki.Core)
	readAll(interpreter, reader, true)
	fmt.Println("done.")

	// repl

	for {

		line, err := rl.Readline()
		if err != nil {
			break
		}

		reader.Append(line)

		if reader.IsBalanced() {
			rl.SetPrompt(promptRepl)
			readAll(interpreter, reader, false)

		} else {
			reader.Append("\n")
			rl.SetPrompt(promptMore)
		}
	}
}

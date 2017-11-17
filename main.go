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
	"github.com/zentrope/embed/interpreter"
)

func eval(lang interpreter.Interpreter, env *interpreter.Environment, form string) interpreter.Expression {
	tokens, err := interpreter.Tokenize(form)
	if err != nil {
		fmt.Printf(" ~ %v\n", err)
		return interpreter.NilExpression
	}

	p := interpreter.NewParser(tokens)
	expr, err := p.Parse()

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return interpreter.NilExpression
	}

	interpreter.Count = 0
	result, err := lang.Evaluate(env, expr)
	fmt.Printf(" - count %v\n", interpreter.Count)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return interpreter.NilExpression
	}

	return result
}

const promptRepl = "repl> "
const promptMore = "   +> "

var flatre = regexp.MustCompile(`\s+`)

func flatten(s string) string {
	return strings.TrimSpace(flatre.ReplaceAllString(s, " "))
}

func readAll(
	lang interpreter.Interpreter,
	env *interpreter.Environment,
	reader *interpreter.Reader, coreload bool) {
	for {
		if reader.IsBalanced() {
			form, err := reader.GetNextForm()
			if err != nil {
				if err == interpreter.ErrEOF {
					break
				}
				fmt.Printf(" ERROR: %v\n", err)
			}
			if form == "" {
				break
			}

			if coreload {
				fmt.Printf("LOADING: %v\n", flatten(form))
				eval(lang, env, form)
			} else {
				fmt.Printf("%v\n", eval(lang, env, form))
			}
			continue
		}
	}
}

func main() {
	fmt.Println("Embed Project Repl")

	rl, err := readline.New(promptRepl)
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	// lang := interpreter.NewInterpreter(interpreter.Naive)
	// fmt.Println("Naive interpreter")
	lang := interpreter.NewInterpreter(interpreter.TCO)
	fmt.Println("TCO interpreter")

	environment := interpreter.NewEnvironment()
	reader := interpreter.NewReader()

	// load core

	reader.Append(interpreter.Core)
	readAll(lang, environment, reader, true)

	// repl

	for {

		line, err := rl.Readline()
		if err != nil {
			break
		}

		reader.Append(line)

		if reader.IsBalanced() {
			rl.SetPrompt(promptRepl)
			readAll(lang, environment, reader, false)

		} else {
			reader.Append("\n")
			rl.SetPrompt(promptMore)
		}
	}
}

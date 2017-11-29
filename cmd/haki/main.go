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

func eval(interpreter haki.Interpreter, env *haki.Environment, form string) haki.Expression {
	tokens, err := haki.Tokenize(form)
	if err != nil {
		fmt.Printf(" ~ %v\n", err)
		return haki.NilExpression
	}

	p := haki.NewParser(tokens)
	expr, err := p.Parse()

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return haki.NilExpression
	}

	result, err := interpreter.Evaluate(env, expr)

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

func readAll(
	interpreter haki.Interpreter,
	env *haki.Environment,
	reader *haki.Reader, coreload bool) {
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
				fmt.Printf("LOADING: %v\n", flatten(form))
				eval(interpreter, env, form)
			} else {
				fmt.Printf("%v\n", eval(interpreter, env, form))
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

	// lang := haki.NewInterpreter(haki.Naive)
	// fmt.Println("Naive interpreter")
	lang := haki.NewInterpreter(haki.TCO)
	fmt.Println("TCO interpreter")

	environment := haki.NewEnvironment()
	reader := haki.NewReader()

	// load core

	reader.Append(haki.Core)
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

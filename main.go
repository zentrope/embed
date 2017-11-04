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

	"github.com/chzyer/readline"
	"github.com/zentrope/embed/interpreter"
)

func exec(tokens *interpreter.Tokens) {
	for i, t := range tokens.Tokens {
		fmt.Printf(" - %2v %v\n", i, t)
	}
}

func eval(form string) {
	tokens, err := interpreter.Tokenize(form)
	if err != nil {
		fmt.Printf(" ~ %v\n", err)
	}
	exec(tokens)
}

const promptRepl = "repl> "
const promptMore = "   +> "

func main() {
	fmt.Println("Embed Project Repl")

	rl, err := readline.New(promptRepl)
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	reader := interpreter.NewReader()

	for {

		line, err := rl.Readline()
		if err != nil {
			break
		}

		reader.Append(line)

		if reader.IsBalanced() {
			rl.SetPrompt(promptRepl)
			form, err := reader.GetNextForm()
			if err != nil {
				fmt.Printf(" ERROR: %v\n", err)
				continue
			}

			eval(form)

		} else {
			reader.Append("\n")
			rl.SetPrompt(promptMore)
		}
	}
}

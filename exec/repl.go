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

package exec

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/zentrope/haki/lang"
)

const promptRepl = "haki> "
const promptMore = "   +> "

// InvokeRepl starts the REPL mode for Haki
func InvokeRepl() {
	printf("Haki Repl")

	rl, err := readline.New(promptRepl)
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	interpreter := lang.NewInterpreter(lang.TCO)
	reader := lang.NewReader(lang.Core)

	// load core

	fmt.Print("* loading core")
	readAll(interpreter, reader, true)
	printf("done.")

	printf("* type :quit to exit")

	// repl

	for {

		line, err := rl.Readline()
		if err != nil {
			printf("bye: %v", err)
			os.Exit(0)
		}

		if line == ":quit" {
			printf("bye")
			os.Exit(0)
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

func printf(pattern string, args ...interface{}) {
	fmt.Printf(pattern+"\n", args...)
}

func eval(interpreter lang.Interpreter, form string) lang.Expression {
	result, err := interpreter.Execute(form)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return lang.NilExpression
	}
	return result
}

func readAll(interpreter lang.Interpreter, reader *lang.Reader, coreload bool) {
	for {
		if reader.IsBalanced() {
			form, err := reader.GetNextForm()
			if err != nil {
				if err == lang.ErrEOF {
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

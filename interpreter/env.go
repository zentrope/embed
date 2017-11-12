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

// Frame contains a single scope of an environment.
type Frame struct {
	data map[string]Sexp
}

// Environment represents the values an expression can see.
type Environment struct {
	frames []Frame
}

// NewEnvironment returns a new environment.
func NewEnvironment() *Environment {
	frame := Frame{
		data: make(map[string]Sexp, 0),
	}

	env := &Environment{
		frames: []Frame{frame},
	}

	for name, fn := range builtins {
		frame.set(name, sexpPrimitive(fn))
	}

	return env
}

// Extend the environment with a new binding.
func (env *Environment) Extend(key string, value Sexp) {
	index := len(env.frames) - 1
	frame := env.frames[index]
	frame.set(key, value)
}

// Lookup a value in the environment.
func (env *Environment) Lookup(key string) (bool, Sexp) {
	high := len(env.frames) - 1
	for i := high; i >= 0; i-- {
		if found, value := env.frames[i].lookup(key); found == true {
			return true, value
		}
	}
	return false, nil
}

func (frame *Frame) lookup(key string) (bool, Sexp) {
	value := frame.data[key]
	if value == nil {
		return false, nil
	}
	return true, value
}

func (frame *Frame) set(key string, value Sexp) {
	frame.data[key] = value
}

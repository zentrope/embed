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

import "fmt"

type frameType map[string]Expression

// Environment represents bindings
type Environment struct {
	global frameType
	frames []frameType
}

// NewEnvironment contains bindings
func NewEnvironment() *Environment {
	data := make(map[string]Expression, 0)
	frames := make([]frameType, 0)

	for name, fn := range builtins {
		data[name] = NewExpr(ExpPrimitive, fn)
	}

	data["true"] = TrueExpression
	data["false"] = FalseExpression
	data["nil"] = NilExpression
	data["&stdin"] = StdinExpression
	data["&stdout"] = StdoutExpression
	data["&stderr"] = StderrExpression

	return &Environment{global: data, frames: frames}
}

// Lookup a value in the environment
func (env *Environment) Lookup(key string) (bool, Expression) {
	for i := len(env.frames) - 1; i >= 0; i-- {
		value, found := env.frames[i].lookup(key)
		if found {
			return true, value
		}
	}

	value, found := env.global.lookup(key)
	if found {
		return true, value
	}
	return false, NilExpression
}

// Set a value in the global environment frame
func (env *Environment) Set(key, value Expression) {
	env.global[key.symbol] = value
}

// Clone returns a copy of the environment
func (env *Environment) Clone() *Environment {
	frames := copyFrames(env.frames)
	return &Environment{
		frames: frames,
		global: env.global,
	}
}

// Dump stack frames
func (env *Environment) Dump() {
	for i, f := range env.frames {
		for k, v := range f {
			fmt.Printf(" frame[%v]: `%v` → `%v`\n", i, k, v)
			if v.IsLambda() {
				v.functionEnv.Dump()
			}
		}
	}
}

// ExtendEnvironment returns a copy of the environment with new bindings.
func (env *Environment) ExtendEnvironment(params Expression, args []Expression) *Environment {
	clone := env.Clone()

	frame := make(frameType, 0)
	for i := 0; i < len(args); i++ {
		frame[params.list[i].symbol] = args[i]
	}

	clone.frames = append(clone.frames, frame)
	return clone
}

func (frame frameType) lookup(key string) (Expression, bool) {
	value, found := frame[key]
	if !found {
		return NilExpression, false
	}
	return value, true
}

// Replace sets an env binding in the current frame
func (env *Environment) Replace(key, value Expression) {
	if len(env.frames) == 0 {
		env.Set(key, value)
		return
	}

	lastFrame := env.frames[len(env.frames)-1]
	lastFrame[key.symbol] = value
}

func copyFrame(frame frameType) frameType {
	newFrame := make(frameType, 0)
	for k, v := range frame {
		newFrame[k] = v
	}
	return newFrame
}

func copyFrames(frames []frameType) []frameType {
	newFrames := make([]frameType, 0)
	for _, frame := range frames {
		newFrames = append(newFrames, copyFrame(frame))
	}
	return newFrames
}

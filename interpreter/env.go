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

	return &Environment{global: data, frames: frames}
}

func (frame frameType) lookup(key string) (Expression, bool) {
	value, found := frame[key]
	if !found {
		return NilExpression, false
	}
	return value, true
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

// Set a value in the current environment frame
func (env *Environment) Set(key Expression, value Expression) {
	env.global[key.symbol] = value
}

func copyFrame(frame frameType) frameType {
	newFrame := make(frameType, len(frame))
	for k, v := range frame {
		newFrame[k] = v
	}
	return newFrame
}

func copyFrames(frames []frameType) []frameType {
	newFrames := make([]frameType, len(frames))
	for _, frame := range frames {
		newFrames = append(newFrames, copyFrame(frame))
	}
	return newFrames
}

// Clone returns a copy of the environment
func (env *Environment) Clone() *Environment {
	// Might want to keep a ptr to the global so all
	// lambdas can reference forward defined vars.
	global := copyFrame(env.global)
	frames := copyFrames(env.frames)
	return &Environment{
		global: global,
		frames: frames,
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

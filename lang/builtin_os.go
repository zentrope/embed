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
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var osBuiltins = primitivesMap{
	"cd!":         _cdBang,
	"cwd":         _cwd,
	"env":         _env,
	"environment": _environment,
	"exec!":       _execBang,
	"exec!!":      _execBangBang,
	"exit!":       _exitBang,
	"shell!":      _shellBang,
}

func toStringSlice(args []Expression) []string {
	params := make([]string, 0)
	for _, a := range args {
		if a.IsString() {
			params = append(params, a.string)
		} else {
			params = append(params, a.String())
		}
	}
	return params
}

func _cdBang(args []Expression) (Expression, error) {
	if err := typeCheck("(cd! path)", args, ckArity(1), ckString(0)); err != nil {
		return NIL, err
	}

	path := args[0].string

	if err := os.Chdir(path); err != nil {
		return NIL, err
	}

	dir, err := os.Getwd()
	if err != nil {
		return NIL, err
	}
	return hStr(dir), nil
}

func _cwd(args []Expression) (Expression, error) {
	if err := typeCheck("(cwd)", args, ckArity(0)); err != nil {
		return NIL, err
	}

	dir, err := os.Getwd()
	if err != nil {
		return NIL, err
	}
	return hStr(dir), nil
}

func _shellBang(args []Expression) (Expression, error) {
	if err := typeCheck("(shell! cmd args…)", args,
		ckArityAtLeast(1), ckString(0)); err != nil {
		return NIL, err
	}

	cmd := args[0].string
	params := toStringSlice(args[1:])

	proc := exec.Command(cmd, params...)

	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	err := proc.Run()

	if err != nil {
		return hStr(err.Error()), nil
	}
	return NIL, nil
}

func _execBang(args []Expression) (Expression, error) {

	if err := typeCheck("(exec! cmd args…)", args,
		ckArityAtLeast(1), ckString(0)); err != nil {
		return NIL, err
	}

	cmd := args[0].string
	params := toStringSlice(args[1:])

	out, err := exec.Command(cmd, params...).CombinedOutput()
	if err != nil {
		return hLst(FALSE, hStr(err.Error()), hStr(string(out))), nil
	}

	return hLst(TRUE, hStr("0"), hStr(string(out))), nil
}

func _execBangBang(args []Expression) (Expression, error) {
	if err := typeCheck("(exec!! cmd args…)", args,
		ckArityAtLeast(1), ckString(0)); err != nil {
		return NIL, err
	}

	cmd := args[0].string
	params := toStringSlice(args[1:])

	proc := exec.Command(cmd, params...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	proc.Stdout = &outBuf
	proc.Stderr = &errBuf

	err := proc.Run()

	outStr := outBuf.String()
	errStr := errBuf.String()

	m := newHakiMap()
	m.set(hSym("ok"), TRUE)
	m.set(hSym("stderr"), hStr(errStr))
	m.set(hSym("stdout"), hStr(outStr))
	m.set(hSym("exit"), hStr("0"))

	if err != nil {
		m.set(hSym("ok"), FALSE)
		m.set(hSym("exit"), hStr(err.Error()))
	}

	return hMap(m), nil
}

func _env(args []Expression) (Expression, error) {
	if err := typeCheck("(env string)", args,
		ckArityAtLeast(1), ckString(0), ckOptString(1)); err != nil {
		return NIL, err
	}

	name := args[0].string

	result := os.Getenv(name)
	if result == "" {
		if len(args) == 2 {
			return hStr(args[1].string), nil
		}
		return NIL, nil
	}
	return hStr(result), nil
}

func _environment(args []Expression) (Expression, error) {

	if err := typeCheck("(environment)", args, ckArity(0)); err != nil {
		return NIL, err
	}

	env := os.Environ()
	results := newHakiMap()

	for _, e := range env {
		words := strings.Split(e, "=")
		results.set(hStr(words[0]), hStr(words[1]))
	}

	return hMap(results), nil
}

func _exitBang(args []Expression) (Expression, error) {
	if err := typeCheck("(exit! [int])", args, ckArityOneOf(0, 1), ckOptInt(1)); err != nil {
		return NIL, err
	}

	code := 0
	if len(args) == 1 {
		code = int(args[0].integer)
	}

	os.Exit(code)
	return NIL, nil
}

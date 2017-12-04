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
	"io/ioutil"
	"os"
)

var fileioBuiltins = primitivesMap{
	"close":     _close,
	"dir?":      _dirP,
	"exists?":   _existsP,
	"file?":     _fileP,
	"read-file": _readFile,
	"handle?":   _handleP,
	"open":      _open,
}

func _readFile(args []Expression) (Expression, error) {
	sig := "(read-file fpath)"
	specs := []spec{ckArity(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	buffer, err := ioutil.ReadFile(args[0].string)
	if err != nil {
		return NilExpression, err
	}

	str := string(buffer)
	return NewExpr(ExpString, str), nil
}

func _open(args []Expression) (Expression, error) {
	sig := "(open fpath)"
	specs := []spec{ckArity(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	path := args[0].string
	file, err := os.Open(path)
	if err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpFile, file), nil
}

func _close(args []Expression) (Expression, error) {
	sig := "(close fhandle)"
	specs := []spec{ckArity(1), ckHandle(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	fileData := args[0].file
	fileData.isOpen = false
	if err := fileData.file.Close(); err != nil {
		return NilExpression, err
	}

	return NilExpression, nil
}

func _dirP(args []Expression) (Expression, error) {
	sig := "(dir? fpath)"
	specs := []spec{ckArity(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	path := args[0].string

	f, err := os.Open(path)
	if err != nil {
		return FalseExpression, nil
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FalseExpression, nil
	}

	return NewExpr(ExpBool, info.IsDir()), nil
}

func _existsP(args []Expression) (Expression, error) {
	sig := "(exists? fpath)"
	specs := []spec{ckArity(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	path := args[0].string

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return TrueExpression, nil
	}
	return FalseExpression, nil
}

func _fileP(args []Expression) (Expression, error) {
	sig := "(file? fpath)"
	specs := []spec{ckArity(1), ckString(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	path := args[0].string

	f, err := os.Open(path)
	if err != nil {
		return FalseExpression, nil
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FalseExpression, nil
	}

	return NewExpr(ExpBool, !info.IsDir()), nil
}

func _handleP(args []Expression) (Expression, error) {
	sig := "(handle? val)"
	specs := []spec{ckArity(1)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpBool, args[0].tag == ExpFile), nil
}

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
	"close":      _close,
	"directory?": _directoryP,
	"exists?":    _existsP,
	"file?":      _fileP,
	"read-file":  _readFile,
	"handle?":    _handleP,
	"open":       _open,
}

func _readFile(args []Expression) (Expression, error) {
	argc := len(args)
	if argc < 1 {
		return nilExpr("(read-file file-name) requires 1 arg, you provided %v", argc)
	}

	if err := verifyStrings(args); err != nil {
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
	argc := len(args)
	if argc != 1 {
		return nilExpr("(open file-path) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(open file-path) expects file-path to be a string, not %v",
			ExprTypeName(args[0].tag))
	}
	path := args[0].string
	file, err := os.Open(path)
	if err != nil {
		return NilExpression, err
	}

	return NewExpr(ExpFile, file), nil
}

func _close(args []Expression) (Expression, error) {

	argc := len(args)
	if argc != 1 {
		return nilExpr("(close file-handle) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpFile {
		return nilExpr("(close file-handle) expects file-handle to be a file-handle, not %v",
			ExprTypeName(args[0].tag))
	}

	fileData := args[0].file
	fileData.isOpen = false
	if err := fileData.file.Close(); err != nil {
		return NilExpression, err
	}

	return NilExpression, nil
}

func _directoryP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(directory? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(directory? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
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
	argc := len(args)
	if argc != 1 {
		return nilExpr("(exists? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(exists? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
	}

	path := args[0].string

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return TrueExpression, nil
	}
	return FalseExpression, nil
}

func _fileP(args []Expression) (Expression, error) {
	argc := len(args)
	if argc != 1 {
		return nilExpr("(file? file-name) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpString {
		return nilExpr("(file? file-name) ← file-name should be string, not %v",
			ExprTypeName(args[0].tag))
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
	argc := len(args)
	if argc != 1 {
		return nilExpr("(handle? file-handle) requires 1 arg, you provided %v", argc)
	}

	if args[0].tag != ExpFile {
		return nilExpr("(handle? file-handle) ← file-handle should be 'handle', not '%v'",
			ExprTypeName(args[0].tag))
	}

	return TrueExpression, nil
}

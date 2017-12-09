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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Used as the payload for file-handle expressions.
type fileData struct {
	file    *os.File
	isOpen  bool
	path    string
	scanner *bufio.Scanner
}

var fileioBuiltins = primitivesMap{
	"close!":    _close,
	"closed?":   _closedP,
	"dir?":      _dirP,
	"files":     _files,
	"exists?":   _existsP,
	"file?":     _fileP,
	"handle?":   _handleP,
	"open!":     _open,
	"read-file": _readFile,
	"read-line": _readLine,
}

// NewFileHandleExpr returns a new file-handle expression.
func NewFileHandleExpr(file *os.File) Expression {
	data := make([]interface{}, 0)
	data = append(data, ExpFile)
	data = append(data, file.Name())

	path, err := filepath.Abs(file.Name())
	if err != nil {
		path = file.Name()
	}

	fileData := &fileData{
		file:    file,
		isOpen:  true,
		path:    path,
		scanner: bufio.NewScanner(file),
	}

	return Expression{tag: ExpFile, hash: hashIt(data...), file: fileData}
}

//-----------------------------------------------------------------------------
// Implementation
//-----------------------------------------------------------------------------

var _rootRegex = ""

func _files(args []Expression) (Expression, error) {

	if err := typeCheck("(dirs path)", args, ckArityAtLeast(1), ckString(0), ckOptString(1)); err != nil {
		return NilExpression, err
	}

	root := args[0].string
	pattern := "*"
	if len(args) == 2 {
		pattern = args[1].string
	}

	matches := make([]string, 0)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ok, err := filepath.Match(pattern, info.Name())
		if ok {
			matches = append(matches, path)
		}
		if err != nil {
			return err
		}
		return nil
	}

	if err := filepath.Walk(root, walker); err != nil {
		return NilExpression, err
	}

	list := make([]Expression, 0)
	for _, match := range matches {
		list = append(list, NewStringExpr(match))
	}

	return NewListExpr(list), nil
}

func _readLine(args []Expression) (Expression, error) {

	if err := typeCheck("(read-line fhandle)", args, ckArity(1), ckHandle(0)); err != nil {
		return NilExpression, err
	}

	fileData := args[0].file

	if !fileData.isOpen || fileData.scanner == nil {
		return NilExpression,
			fmt.Errorf("Cannot read from un-opened file: '%v'",
				fileData.path)
	}

	moreToScan := fileData.scanner.Scan()

	if !moreToScan {
		err := fileData.scanner.Err()
		if err == nil {
			fileData.file.Close()
			fileData.isOpen = false
			fileData.scanner = nil
			return NilExpression, nil
		}

		return NilExpression, err
	}

	line := fileData.scanner.Text()
	return NewStringExpr(line), nil
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

	return NewFileHandleExpr(file), nil
}

func safeClose(f *os.File) error {
	if err := f.Close(); err != nil {
		if pe, ok := err.(*os.PathError); !ok {
			return err
		} else if pe.Err != os.ErrClosed {
			return err
		}
	}
	return nil
}

func _close(args []Expression) (Expression, error) {
	sig := "(close fhandle)"
	specs := []spec{ckArity(1), ckHandle(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	fileData := args[0].file
	fileData.isOpen = false
	fileData.scanner = nil
	if err := safeClose(fileData.file); err != nil {
		return NilExpression, err
	}
	return NilExpression, nil
}

func _closedP(args []Expression) (Expression, error) {
	sig := "(closed? fhandle)"
	specs := []spec{ckArity(1), ckHandle(0)}

	if err := typeCheck(sig, args, specs...); err != nil {
		return NilExpression, err
	}

	fileData := args[0].file

	return NewBoolExpr(!fileData.isOpen), nil
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

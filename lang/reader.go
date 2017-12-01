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

package lang

import (
	"errors"
	"fmt"
	"strings"
)

// Reader is a container for dolling out expressions.
type Reader struct {
	buffer []rune
}

// NewReader returns a new instances of a reader
func NewReader(forms ...string) *Reader {
	r := &Reader{
		buffer: make([]rune, 0),
	}

	for _, f := range forms {
		r.Append(f)
	}

	return r
}

func (reader *Reader) String() string {
	return fmt.Sprintf("#<reader:%#v>", string(reader.buffer))
}

// IsBalanced returns true if a given string represents a complete
// expression.
func (reader *Reader) IsBalanced() bool {
	opens := 0
	closes := 0
	for _, c := range reader.buffer {
		switch c {
		case '(':
			opens = opens + 1
		case ')':
			closes = closes + 1
		}
	}
	return opens == closes
}

// Append appends new data to the reader.
func (reader *Reader) Append(line string) {
	reader.buffer = append(reader.buffer, []rune(line)...)
}

// ErrEOF means there's nothing left to read
var ErrEOF = errors.New("EOF")

// GetForms returns all the available forms in the reader buffer.
func (reader *Reader) GetForms() ([]string, error) {
	forms := make([]string, 0)

	for {
		form, err := reader.GetNextForm()
		if err == ErrEOF {
			return forms, nil
		} else if err != nil {
			return []string{}, err
		}

		if form == "" {
			break
		}

		forms = append(forms, form)
	}
	return forms, nil
}

// GetNextForm returns the next available expression from the buffer.
func (reader *Reader) GetNextForm() (string, error) {

	if len(reader.buffer) == 0 {
		return "", ErrEOF
	}

	form := make([]rune, 0)

	opens := 0
	closes := 0

	for _, c := range reader.buffer {
		if c == '(' {
			opens = opens + 1
		} else if c == ')' {
			closes = closes + 1
		}

		form = append(form, c)

		if opens > 0 && (opens == closes) {
			break
		}
	}

	reader.buffer = reader.buffer[len(form):]

	if opens != closes {
		return string(form), errors.New("incomplete form (missing parens)")
	}
	return strings.TrimSpace(string(form)), nil
}

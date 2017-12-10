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
	"io/ioutil"
	"regexp"

	"github.com/zentrope/haki/lang"
)

// InvokeScript loads and runs a haki script.
func InvokeScript(filename string, args []string) error {

	script, err := loadScript(filename)
	if err != nil {
		return err
	}

	return runScript(script, args)
}

var hashBangRe = regexp.MustCompile("(?m)^[#][!].*$")

func loadScript(fname string) (string, error) {
	buffer, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}

	str := string(buffer)
	return hashBangRe.ReplaceAllString(str, ""), nil
}

func runScript(script string, args []string) error {
	interpreter := lang.NewScriptInterpreter(lang.TCO, args)
	reader := lang.NewReader(lang.Core, script)

	_, err := interpreter.Run(reader)
	return err
}

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

	"github.com/zentrope/haki/lang"
)

var binVers = "1"
var binCommit = "dev"
var binDate = "dev"

func version() string {
	return fmt.Sprintf("(vers: %#v, commit: %#v, date: %#v)", binVers, binCommit, binDate)
}

func setVersionEnv(haki lang.Interpreter) {
	haki.SetVersionInfo(binVers, binCommit, binDate)
}

/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package cmdmux

import (
	"fmt"
	"io"
	"os"
)

// OutputCompletion outputs the command completion file for bash.
func (c *CmdMux) OutputCompletion(w io.Writer) error {
	return func(w io.Writer, node *cmdNode) error {

		program := os.Args[0]
		fmt.Fprintf(w, "# bash completion file for %s", program)
		fmt.Fprintf(w, "# Copy this file to somewhere (e.g. ~/.test-completion)\n"+
			"# and then `$ source ~/.test-completion`\n")

		fmt.Fprintf(w, "_%s()\n{\n", program)

		fmt.Fprintf(w, "  local cur prev opts\n")
		fmt.Fprintf(w, "  COMPREPLY=()\n")
		fmt.Fprintf(w, "  cur=\"${COMP_WORDS[COMP_CWORD]}\"\n")
		fmt.Fprintf(w, "  cur=\"${COMP_WORDS[COMP_CWORD]}\"\n")
		fmt.Fprintf(w, "  opts=\"")
		for _, v := range std.root.subNodes {
			fmt.Fprintf(w, "%s ", v.name)
		}
		fmt.Fprintf(w, "\"\n\n")

		fmt.Fprintf(w, "}\ncomplete -F _%s %s\n", program, program)

		return nil
	}(w, c.root)
}

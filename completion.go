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
)

func outputCmds(w io.Writer, node *cmdNode) {
}

// OutputCompletion outputs the command completion file for *bash*.
func (c *CmdMux) OutputCompletion(program string, w io.Writer) error {
	node := c.root
	fmt.Fprintf(w, "# bash completion file for %s\n", program)
	fmt.Fprintf(w, "# Copy this file to somewhere (e.g. ~/.test-completion)\n"+
		"# and then `$ source ~/.test-completion`\n\n")

	fmt.Fprintf(w, "_%s()\n{\n", program)

	fmt.Fprintln(w, `  local cur prev opts`)
	fmt.Fprintln(w, `  COMPREPLY=()`)
	fmt.Fprintln(w, `  cur="${COMP_WORDS[COMP_CWORD]}"`)
	fmt.Fprintln(w, `  prev="${COMP_WORDS[COMP_CWORD-1]}"`)
	fmt.Fprintf(w, `  opts="`)
	for _, v := range node.subNodes {
		fmt.Fprintf(w, "%s ", v.name)
	}
	fmt.Fprintf(w, "\"\n\n")
	fmt.Fprintln(w, `  case "$prev" in`)

	outputCmds(w, node)

	fmt.Fprintln(w, "  *)")
	fmt.Fprintln(w, `    local prev2="${COMP_WORDS[COMP_CWORD-2]}"`)
	fmt.Fprintf(w, "    ;;\n")
	fmt.Fprintf(w, "  esac\n\n")
	fmt.Fprintln(w, `  COMPREPLY=( $(compgen -W "$opts" -- $cur) )`)
	fmt.Fprintln(w, `  return 0`)
	fmt.Fprintf(w, "}\ncomplete -F _%s %s\n", program, program)

	return nil
}

func OutputCompletion(program string, w io.Writer) error {
	return std.OutputCompletion(program, w)
}

// _kbdashboard()
// {
//     local cur prev opts
//     COMPREPLY=()
//     cur="${COMP_WORDS[COMP_CWORD]}"
//     prev="${COMP_WORDS[COMP_CWORD-1]}"
//     opts="list choose edit help make build config install version"
//
//     case "$prev" in
//     choose | make | install | version)
//         COMPREPLY=()
//         return 0
//         ;;
//     list )
//         COMPREPLY=( $(compgen -W "-v" -- $cur) )
//         return 0
//         ;;
//     edit )
//         COMPREPLY=( $(compgen -W "profile install" -- $cur) )
//         return 0
//         ;;
//     *)
//         local prev2="${COMP_WORDS[COMP_CWORD-2]}"
//         ;;
//     esac
//
//     COMPREPLY=( $(compgen -W "$opts" -- $cur) )
//     return 0
// }
// complete -F _kbdashboard kbdashboard

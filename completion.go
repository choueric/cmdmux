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

const (
	header = `# Copy this file to somewhere (e.g. ~/.test-completion)
# and then '$ source ~/.test-completion')

`
	body1 = `  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
`
	end = `  *)
    local prev2="${COMP_WORDS[COMP_CWORD-2]}"
    ;;
  esac

  COMPREPLY=( $(compgen -W "$opts" -- $cur) )
  return 0
}
`
	entryHead = `    COMPREPLY=( $(compgen -W "`
	entryEnd  = `" -- $cur) )
    return 0
    ;;
`
)

func generateEntry(node *cmdNode, depth int, data interface{}) {
	if len(node.subNodes) == 0 {
		return
	}
	w := data.(io.Writer)
	fmt.Fprintf(w, "  %s )\n", node.name)
	fmt.Fprintf(w, entryHead)
	for _, sub := range node.subNodes {
		fmt.Fprintf(w, "%s ", sub.name)
	}
	fmt.Fprintf(w, entryEnd)
}

func generateEntries(w io.Writer, node *cmdNode) {
	walkByDepth(node.subNodes, 0, generateEntry, w)
}

// GenerateCompletion generate shell completion file for *bash*.
//
// 1. list all depth 1 nodes into $opts
// 2. create entry for every node which has sub-node.
func (c *CmdMux) GenerateCompletion(program string, w io.Writer) error {
	fmt.Fprintf(w, "# bash completion file for %s\n", program)
	fmt.Fprintf(w, header)

	fmt.Fprintf(w, "_%s()\n{\n", program)
	fmt.Fprintf(w, body1)

	// 1. create opts
	fmt.Fprintf(w, `  opts="`)
	for _, v := range c.root.subNodes {
		fmt.Fprintf(w, "%s ", v.name)
	}
	fmt.Fprintf(w, "\"\n\n")
	fmt.Fprintln(w, `  case "$prev" in`)

	// 2. create entries
	generateEntries(w, c.root)

	fmt.Fprintf(w, end)
	fmt.Fprintf(w, "complete -F _%s %s\n", program, program)

	return nil
}

func GenerateCompletion(program string, w io.Writer) error {
	return std.GenerateCompletion(program, w)
}

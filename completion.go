package cmdmux

import (
	"fmt"
	"io"
)

const (
	header = "# bash completion file for %s\n" +
		"# Copy this file to somewhere (e.g. ~/.%s-completion)\n" +
		"# and then '$ source ~/.%s-completion'\n\n"

	body = `  local cur prev opts
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

// GenerateCompletion generates a *bash* completion file with the program name.
func (c *CmdMux) GenerateCompletion(program string, w io.Writer) error {
	fmt.Fprintf(w, header, program, program, program)

	fmt.Fprintf(w, "_%s()\n{\n", program)
	fmt.Fprintf(w, body)

	// 1. list all depth 1 nodes into $opts
	fmt.Fprintf(w, `  opts="`)
	for _, v := range c.root.subNodes {
		fmt.Fprintf(w, "%s ", v.name)
	}
	fmt.Fprintf(w, "\"\n\n")
	fmt.Fprintln(w, `  case "$prev" in`)

	// 2. create entry for every node which has sub-node.
	generateEntries(w, c.root)

	fmt.Fprintf(w, end)
	fmt.Fprintf(w, "complete -F _%s %s\n", program, program)

	return nil
}

// GenerateCompletion generates a *bash* completion file with the program name.
func GenerateCompletion(program string, w io.Writer) error {
	return std.GenerateCompletion(program, w)
}

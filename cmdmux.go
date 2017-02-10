// Package cmdmux is used to parse and route commands of terminal program.
package cmdmux

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// CmdHanlder is the type of callback function for command.
// if error is nil, then return value int is useful.
type CmdHandler func([]string, interface{}) (int, error)

// CmdMux represents program's commands.
type CmdMux struct {
	root *cmdNode
}

var std = New()

// New creates a CmdMux with the default "/" empty command.
func New() *CmdMux {
	n := newCmdNode("/")
	c := &CmdMux{root: n}
	return c
}

// PrintTree outputs a simple tree structure of c
func (c *CmdMux) PrintTree(w io.Writer) {
	c.root.printTree(w)
}

// HandleFunc registers the handler function for the given command path cmdpath
func (c *CmdMux) HandleFunc(cmdpath string, handler CmdHandler) error {
	node, err := c.newNode(cmdpath)
	if err != nil {
		return err
	}
	node.ops.Handler = handler
	//fmt.Printf("cmdmux: add handler %s -> %s\n", cmdpath, node.name)

	return nil
}

// AddHelpInfo adds help information to the registed cmd node.
func (c *CmdMux) AddHelpInfo(cmdpath string, synopsis func() string, usage func() string) error {
	node, err := c.getNode(cmdpath)
	if err != nil {
		return err
	}
	node.ops.Synopsis = synopsis
	node.ops.Usage = usage
	return nil
}

// Execute accepts the os.Args as command and executes it with data
func (c *CmdMux) Execute(data interface{}) (int, error) {
	path := ""
	node := c.root
	var opts []string
	args := os.Args[1:]
	for i, v := range args {
		sub := node.hasSubNode(v)
		if sub == nil {
			opts = args[i:]
			break
		}
		node = sub
		path = path + "/" + node.name
	}

	if path == "" {
		path = "/"
	}
	if node.ops.Handler == nil {
		msg := fmt.Sprintf("cmdmux: %s does not have a handler.", path)
		c.PrintUsage(os.Stderr, "")
		return 0, errors.New(msg)
	}

	//fmt.Printf("cmdmux: invode %s\n", path)
	return node.ops.Handler(opts, data)
}

// PrintUsage outputs usage of cmd node specified by cmdpath.
// cmdpath must be in the format like "/cmd1/cmd2".
func (c *CmdMux) PrintUsage(w io.Writer, cmdpath string) {
	if cmdpath == "" {
		c.root.printAllUsages(w, 0)
		return
	}

	node, err := c.getNode(cmdpath)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	if node.ops.Usage != nil {
		fmt.Fprintf(os.Stderr, "%s", node.ops.Usage())
	}
}

// HandleFunc registers the handler function for the given command path cmdpath
// in the default CmdMux.
func HandleFunc(cmdpath string, handler CmdHandler) error {
	return std.HandleFunc(cmdpath, handler)
}

// AddHelpInfo adds help information to the registed cmd node
// in the default CmdMux.
func AddHelpInfo(cmdpath string, synopsis func() string, usage func() string) error {
	return std.AddHelpInfo(cmdpath, synopsis, usage)
}

// Execute accepts the os.Args as command and executes it with data
// in the default CmdMux
func Execute(data interface{}) (int, error) {
	return std.Execute(data)
}

// PrintTree outputs a simple tree structure of built-in cmdmux
func PrintTree(w io.Writer) {
	std.PrintTree(w)
}

func PrintUsage(w io.Writer, cmdpath string) {
	std.PrintUsage(w, cmdpath)
}

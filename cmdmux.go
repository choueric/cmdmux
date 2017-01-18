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

// Package cmdmux is used to parse and route commands of terminal programe.
package cmdmux

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// CmdHanlder is the type of callback function for command.
// if error is nil, then return value int is useful.
type CmdHandler func([]string, interface{}) (int, error)

// CmdMux represents programme's commands.
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

// String return the string format of CmdMux c.
func (c *CmdMux) String() string {
	var result string
	c.root.toString("/", &result)
	return result
}

// PrintTree outputs a simple tree structure of c
func (c *CmdMux) PrintTree(w io.Writer) {
	c.root.printTree(w)
}

// HandleFunc registers the handler function for the given command path cmdpath
func (c *CmdMux) HandleFunc(cmdpath string, handler CmdHandler) error {
	if cmdpath[0] != '/' {
		return errors.New("cmdmux: cmdpath should be absolute")
	}

	if cmdpath == "/" {
		c.root.handler = handler
		return nil
	}

	cmdStrs := strings.Split(cmdpath, "/")[1:]
	last := len(cmdStrs) - 1
	node := c.root
	for i, v := range cmdStrs {
		sub := node.hasSubNode(v)
		if sub == nil {
			sub = newCmdNode(v)
			node.subNodes = append(node.subNodes, sub)
		}
		if i == last {
			node = sub
			break
		}
		node = sub
	}
	node.handler = handler
	//fmt.Printf("cmdmux: add handler %s -> %s\n", cmdpath, node.name)

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
	if node.handler == nil {
		msg := fmt.Sprintf("cmdmux: %s does not have a handler.", path)
		return 0, errors.New(msg)
	}

	//fmt.Printf("cmdmux: invode %s\n", path)
	return node.handler(opts, data)
}

// HandleFunc registers the handler function for the given command path cmdpath
// in the default CmdMux.
func HandleFunc(cmdpath string, handler CmdHandler) error {
	return std.HandleFunc(cmdpath, handler)
}

// Execute accepts the os.Args as command and executes it with data
// in the default CmdMux
func Execute(data interface{}) (int, error) {
	return std.Execute(data)
}

// String() return the string format of default CmdMux
func String() string {
	return std.String()
}

// PrintTree outputs a simple tree structure of built-in cmdmux
func PrintTree(w io.Writer) {
	std.PrintTree(w)
}

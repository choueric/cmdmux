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

type cmdNode struct {
	name     string
	subNodes []*cmdNode
	handler  CmdHandler
}

var std = New()

// New creates a CmdMux variable with the default "/" empty command.
func New() *CmdMux {
	n := &cmdNode{name: "/"}
	c := &CmdMux{root: n}
	return c
}

// String() return the string format of CmdMux c.
func (c *CmdMux) String() string {
	return func(n *cmdNode) string {
		var result string
		n.toString("", &result)
		return result
	}(c.root)
}

// HandleFunc registers the handler function for the given command path cmdpath
func (c *CmdMux) HandleFunc(cmdpath string, handler CmdHandler) error {
	return func(n *cmdNode, cmdpath string, handler CmdHandler) error {
		if cmdpath[0] != '/' {
			return errors.New("cmdmux: cmdpath should be absolute")
		}

		if cmdpath == "/" {
			n.handler = handler
			return nil
		}

		cmdStrs := strings.Split(cmdpath, "/")[1:]
		last := len(cmdStrs) - 1
		node := n
		for i, v := range cmdStrs {
			sub := node.hasSubNode(v)
			if sub == nil {
				sub = &cmdNode{name: v}
				if i == last {
					sub.handler = handler
				}
				node.subNodes = append(node.subNodes, sub)
			}
			node = sub
		}

		return nil
	}(c.root, cmdpath, handler)
}

// Execute accepts the os.Args as command and executes it with data
func (c *CmdMux) Execute(data interface{}) (int, error) {
	return func(node *cmdNode, args []string, data interface{}) (int, error) {
		var opts []string
		args = args[1:]
		for i, v := range args {
			sub := node.hasSubNode(v)
			if sub == nil {
				opts = args[i:]
				break
			}
			node = sub
		}

		if node == nil {
			return 0, errors.New("cmdmux: no such cmdNode")
		}

		if node.handler == nil {
			return 0, errors.New("cmdmux: cmdNode does not have a handler")
		}

		return node.handler(opts, data)
	}(c.root, os.Args, data)
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

func (n *cmdNode) hasSubNode(name string) *cmdNode {
	for _, v := range n.subNodes {
		if v.name == name {
			return v
		}
	}
	return nil
}

func (n *cmdNode) toString(prefix string, result *string) {
	switch prefix {
	case "/":
		prefix = prefix + n.name
	case "":
		prefix = "/"
	default:
		prefix = prefix + "/" + n.name
	}

	if len(n.subNodes) == 0 {
		*result = *result + prefix + "\n"
	} else {
		for _, v := range n.subNodes {
			v.toString(prefix, result)
		}
	}
}

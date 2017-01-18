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
	"errors"
	"fmt"
	"strings"
)

// the reason of using tree instead of map[string]Handler is that:
// 1. it hard to tell the cmd-path from parameters if use map. for example:
//      add("/build", nil)
//      add("/build/kernel", nil)
//      $ test build kernel image
//    here 'image' is actually the parameter, but you don't know if the
//    cmd-path is build, build/kernel or build/kernel/image
//    But use tree, searching from top to bottom, it can definitly find the
//    leaf node.
// 2. it is easy to generate the completion file using tree.
type cmdNode struct {
	name     string
	subNodes []*cmdNode
	handler  CmdHandler
}

type walkHandler func(*cmdNode, int, interface{})

func newCmdNode(name string) *cmdNode {
	node := &cmdNode{name: name}
	return node
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

func (n *cmdNode) depth(preDepth int) int {
	maxDepth := preDepth
	if len(n.subNodes) != 0 {
		for _, v := range n.subNodes {
			d := v.depth(preDepth + 1)
			if maxDepth < d {
				maxDepth = d
			}
		}
	}

	return maxDepth
}

func walkByDepth(nodes []*cmdNode, depth int, f walkHandler, data interface{}) {
	if len(nodes) == 0 {
		return
	}
	depth = depth + 1
	var next []*cmdNode
	for _, n := range nodes {
		f(n, depth, data)
		for _, sub := range n.subNodes {
			next = append(next, sub)
		}
	}
	walkByDepth(next, depth, f, data)
}

func (c *CmdMux) getNode(cmdpath string) (*cmdNode, error) {
	if cmdpath[0] != '/' {
		return nil, errors.New("cmdmux: cmdpath should be absolute")
	}

	if cmdpath == "/" {
		return c.root, nil
	}

	cmdStrs := strings.Split(cmdpath, "/")[1:]
	node := c.root
	for _, v := range cmdStrs {
		sub := node.hasSubNode(v)
		if sub == nil {
			return nil, errors.New(fmt.Sprintf("cmdmux: node %s does not exist.", v))
		}
		node = sub
	}

	return node, nil
}

func (n *cmdNode) hasSubNode(name string) *cmdNode {
	for _, v := range n.subNodes {
		if v.name == name {
			return v
		}
	}
	return nil
}

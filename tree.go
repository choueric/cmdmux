package cmdmux

import (
	"errors"
	"fmt"
	"strings"
)

type cmdNode struct {
	name     string
	subNodes []*cmdNode
	handler  CmdHandler
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

func (c *CmdMux) getCmdNode(cmdpath string) (*cmdNode, error) {
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

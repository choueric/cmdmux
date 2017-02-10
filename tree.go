package cmdmux

import (
	"errors"
	"fmt"
	"io"
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
	ops      CmdOps
	subNodes []*cmdNode
}

// cmdnode, depth, data
type walkHandler func(*cmdNode, int, interface{})

func newCmdNode(name string) *cmdNode {
	node := &cmdNode{name: name}
	return node
}

const exeSymbol = "*"

func (n *cmdNode) modifyName() string {
	if n.ops.Handler == nil {
		return n.name
	} else {
		return n.name + exeSymbol
	}
}

func (n *cmdNode) doPrintTree(w io.Writer, depth int, last bool, onlyOne []bool) {
	for i := 0; i < depth-1; i++ {
		if onlyOne[i] {
			fmt.Fprintf(w, "    ")
		} else {
			fmt.Fprintf(w, "│   ")
		}
	}
	if depth != 0 {
		if last {
			fmt.Fprintln(w, "└── "+n.modifyName())
		} else {
			fmt.Fprintln(w, "├── "+n.modifyName())
		}
	} else {
		fmt.Fprintln(w, n.modifyName())
	}
	len := len(n.subNodes)
	if len == 1 {
		onlyOne = append(onlyOne, true)
	} else {
		onlyOne = append(onlyOne, false)
	}
	for i, subNode := range n.subNodes {
		last := false
		if i == len-1 {
			last = true
		}
		subNode.doPrintTree(w, depth+1, last, onlyOne)
	}
}

func (n *cmdNode) printTree(w io.Writer) {
	var onlyOne []bool
	n.doPrintTree(w, 0, false, onlyOne)
}

func (n *cmdNode) printAllUsages(w io.Writer, depth int) {
	for i := 0; i < depth-1; i++ {
		fmt.Fprintf(w, "  ")
	}
	if n.name != "/" {
		fmt.Fprintf(w, "  %s", n.name)
		if n.ops.Synopsis != nil {
			fmt.Fprintf(w, "\t: %s\n", n.ops.Synopsis())
		} else {
			fmt.Fprintf(w, "\n")
		}
		// TODO: need print usage ?
		/*
			if n.ops.usage != nil {
				n.ops.usage()
			}
		*/
	}

	for _, subNode := range n.subNodes {
		subNode.printAllUsages(w, depth+1)
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

func (c *CmdMux) newNode(cmdpath string) (*cmdNode, error) {
	if cmdpath[0] != '/' {
		return nil, errors.New("cmdmux: cmdpath should be absolute")
	}

	if cmdpath == "/" {
		return c.root, nil
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
	return node, nil
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

////////////////////////////////////////////////////////////////////////////////
var lastDepth = 0

func printLevelNode(node *cmdNode, depth int, data interface{}) {
	w := data.(io.Writer)
	if lastDepth != depth {
		fmt.Fprintf(w, "\n------------ [%d] -----------------\n", depth)
		lastDepth = depth
	}
	if node.ops.Handler != nil {
		fmt.Fprintf(w, "%s ", node.name+exeSymbol)
	} else {
		fmt.Fprintf(w, "%s ", node.name)
	}
}

func (n *cmdNode) printLevels(w io.Writer) {
	walkByDepth(n.subNodes, 0, printLevelNode, w)
	fmt.Fprintf(w, "\n-----------------------------------\n")
}

func (n *cmdNode) toString(prefix string, result *string) {
	switch prefix {
	case "/":
		if n.ops.Handler != nil {
			prefix = exeSymbol
		} else {
			prefix = ""
		}
	default:
		if n.ops.Handler != nil {
			prefix = prefix + "/" + n.name + exeSymbol
		} else {
			prefix = prefix + "/" + n.name
		}
	}

	if len(n.subNodes) == 0 {
		*result = *result + prefix + "\n"
	} else {
		for _, v := range n.subNodes {
			v.toString(prefix, result)
		}
	}
}

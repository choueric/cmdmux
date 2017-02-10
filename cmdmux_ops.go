package cmdmux

import (
	"fmt"
	"os"
)

type CmdOps struct {
	Synopsis func() string
	Usage    func() string
	Handler  CmdHandler
}

func helpUsage() string {
	return "Print help message.\n"
}

func (c *CmdMux) helpHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		c.PrintUsage(os.Stderr, "")
	} else {
		c.PrintUsage(os.Stderr, args[0])
	}
	return 0, nil
}

func (c *CmdMux) Register(cmdpath string, ops CmdOps) error {
	node, err := c.newNode(cmdpath)
	if err != nil {
		return err
	}
	node.ops = ops

	return nil
}

func (c *CmdMux) EnableHelp() {
	c.HandleFunc("/help", c.helpHandler)
	c.AddHelpInfo("/help", func() string { return "show help message" }, helpUsage)
}

func (c *CmdMux) DupOps(destPath, srcPath string) error {
	srcNode, err := c.getNode(srcPath)
	if err != nil {
		return err
	}

	destNode, err := c.getNode(destPath)
	if err != nil {
		return err
	}

	destNode.ops = srcNode.ops

	return nil
}

func Register(cmdpath string, ops CmdOps) error {
	return std.Register(cmdpath, ops)
}

func EnableHelp() {
	std.EnableHelp()
}

func DupOps(destPath, srcPath string) error {
	return std.DupOps(destPath, srcPath)
}

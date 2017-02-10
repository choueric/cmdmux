package main

import (
	"fmt"
	"os"

	"github.com/choueric/cmdmux"
)

var buildKernelImageOps = cmdmux.CmdOps{
	Synopsis: func() string { return "build linux kernel image" },
	Usage:    func() string { return "--> build kernel image usage" },
	Handler:  buildKernelImageHandler,
}

func buildKernelImageHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'build kernel image', args = %v\n", args)
	return 2, nil
}

func buildKernelDtbUsage() string {
	return "--> build kernel dtb usage"
}

func buildKernelDtbHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'build kernel dtb', args = %v\n", args)
	return 2, nil
}

func completionUsage() string {
	return "--> completion usage"
}

func completionHandler(args []string, data interface{}) (int, error) {
	cmdmux.GenerateCompletion("example", os.Stdout)
	return 2, nil
}

func main() {
	cmdpath := "/build/kernel/image"
	cmdmux.Register(cmdpath, buildKernelImageOps)

	cmdpath = "/build/kernel/dtb"
	cmdmux.HandleFunc(cmdpath, buildKernelDtbHandler)
	cmdmux.AddHelpInfo(cmdpath, func() string { return "build liux DTB file" }, buildKernelDtbUsage)

	cmdpath = "/completion"
	cmdmux.HandleFunc(cmdpath, completionHandler)
	cmdmux.AddHelpInfo(cmdpath, func() string { return "generate compeltion file" }, completionUsage)

	cmdmux.EnableHelp()

	ret, err := cmdmux.Execute(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nreturn value: %d\n", ret)
}

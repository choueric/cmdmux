package main

import (
	"fmt"
	"os"

	"github.com/choueric/cmdmux"
)

type Options struct {
	name string
}

func rootHandler(args []string, data interface{}) (int, error) {
	fmt.Println("Usage: build [uboot|kernel]")
	fmt.Println("       build kernel [image|dtb]")
	fmt.Println("       completion")
	return 0, nil
}

func buildHandler(args []string, data interface{}) (int, error) {
	opt := data.(*Options)
	fmt.Printf("invoke 'build' of %s\n", opt.name)
	return 1, nil
}

func buildKernelHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'build kernel', args = %v\n", args)
	return 2, nil
}

func completionHandler(args []string, data interface{}) (int, error) {
	cmdmux.GenerateCompletion("example", os.Stdout)
	return 2, nil
}

func main() {
	opt := &Options{name: "arm"}

	cmdmux.HandleFunc("/", rootHandler)
	cmdmux.HandleFunc("/build/uboot", buildKernelHandler)
	cmdmux.HandleFunc("/build/kernel/image", buildKernelHandler)
	cmdmux.HandleFunc("/build/kernel/dtb", buildKernelHandler)
	cmdmux.HandleFunc("/completion", completionHandler)

	cmdmux.Execute(opt)
}

package main

import (
	"fmt"

	"github.com/choueric/cmdmux"
)

type Options struct {
	name string
}

func rootHandler(args []string, data interface{}) (int, error) {
	fmt.Println("Usage: build")
	fmt.Println("       build kernel")
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

func main() {
	opt := &Options{name: "arm"}

	cmdmux.HandleFunc("/", rootHandler)
	cmdmux.HandleFunc("/build", buildHandler)
	cmdmux.HandleFunc("/build/kernel", buildKernelHandler)

	cmdmux.Execute(opt)
}

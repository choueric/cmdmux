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
	fmt.Println("Usage:")
	cmdmux.PrintTree(os.Stdout)
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

func optionOHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'option -o', args = %v\n", args)
	return 3, nil
}

func optionPHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'option -p', args = %v\n", args)
	return 3, nil
}

func main() {
	opt := &Options{name: "arm"}

	cmdmux.HandleFunc("/", rootHandler)
	cmdmux.HandleFunc("/build/kernel/image", buildKernelHandler)
	cmdmux.HandleFunc("/build/kernel/dtb", buildKernelHandler)
	cmdmux.HandleFunc("/build", buildHandler)
	cmdmux.HandleFunc("/option/-o", optionOHandler)
	cmdmux.HandleFunc("/option/-p", optionPHandler)
	cmdmux.HandleFunc("/completion", completionHandler)

	ret, err := cmdmux.Execute(opt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("return value: %d\n", ret)
}

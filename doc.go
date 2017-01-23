// Simple Example, used like package `http`:
//
//  package main
//
//  import (
//  	"fmt"
//  	"os"
//
//  	"github.com/choueric/cmdmux"
//  )
//
//  type Options struct {
//  	arch string
//  }
//
//  func rootHandler(args []string, data interface{}) (int, error) {
//  	fmt.Println("Usage:")
//  	cmdmux.PrintTree(os.Stderr)
//  	return 0, nil
//  }
//
//  func buildHandler(args []string, data interface{}) (int, error) {
//  	opt := data.(*Options)
//  	fmt.Printf("invoke 'build' of %s\n", opt.arch)
//  	return 1, nil
//  }
//
//  func buildKernelHandler(args []string, data interface{}) (int, error) {
//  	fmt.Printf("invoke 'build kernel', args = %v\n", args)
//  	return 2, nil
//  }
//
//  func main() {
//  	opt := &Options{arch: "arm"}
//
//  	cmdmux.HandleFunc("/", rootHandler)
//  	cmdmux.HandleFunc("/build", buildHandler)
//  	cmdmux.HandleFunc("/build/kernel", buildKernelHandler)
//  	cmdmux.HandleFunc("/build/kernel/image", buildKernelHandler)
//  	cmdmux.HandleFunc("/build/uboot", buildKernelHandler)
//
//  	cmdmux.Execute(opt)
//  }
package cmdmux

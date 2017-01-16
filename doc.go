// Simple Example, used like package `http`:
//  package main
//
//  import (
//  	"fmt"
//
//  	"github.com/choueric/cmdmux"
//  )
//
//  func rootHandler(args []string, data interface{}) (int, error) {
//  	fmt.Println("Usage: build")
//  	fmt.Println("       build kernel")
//  	return 0, nil
//  }
//
//  func buildHandler(args []string, data interface{}) (int, error) {
//  	fmt.Println("invoke 'build'")
//  	return 1, nil
//  }
//
//  func buildKernelHandler(args []string, data interface{}) (int, error) {
//  	fmt.Println("invoke 'build kernel'")
//  	return 2, nil
//  }
//
//  func main() {
//  	cmdmux.HandleFunc("/", rootHandler)
//  	cmdmux.HandleFunc("/build", buildHandler)
//  	cmdmux.HandleFunc("/build/kernel", buildKernelHandler)
//
//  	cmdmux.Execute(nil)
//  }
package cmdmux

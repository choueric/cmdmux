// Simple Example, used like package `http`:
//
//  package main
//
//  import (
//  	"fmt"
//
//  	"github.com/choueric/cmdmux"
//  )
//
//  type Options struct {
//  	profile string
//  }
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
//  	options := data.(*Options)
//  	if options.profile == "" {
//  		fmt.Println("invoke 'build kernel'")
//  	} else {
//  		fmt.Printf("invoke 'build kernel' for %s", options.profile)
//  	}
//  	return 2, nil
//  }
//
//  func main() {
//  	var options Options
//
//  	cmdmux.HandleFunc("/", rootHandler)
//  	cmdmux.HandleFunc("/build", buildHandler)
//  	cmdmux.HandleFunc("/build/kernel", buildKernelHandler)
//
//  	if flagSet, err := cmdmux.FlagSet("/build/kernel"); err == nil {
//  		flagSet.StringVar(&options.profile, "p", "", "speicify profile")
//  	}
//
//  	cmdmux.Execute(&options)
//  }
package cmdmux

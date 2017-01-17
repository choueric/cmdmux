# cmdmux

Package cmdmux implements a command parser and router for terminal programme.

# Overview

In general, there are two styles a terminal programme to interact with users.

1. Use -o, -p to specify parameters. Most programmes are in this way.
2. Use sub-commands, like git which uses only one level sub-command.

The first way can be implemented by the `flag` package of Golang, and the
second this package.

# Usage

Below is a simple example:

```go
package main

import (
	"fmt"

	"github.com/choueric/cmdmux"
)

type Options struct {
	arch string
}

func rootHandler(args []string, data interface{}) (int, error) {
	fmt.Println("Usage: build")
	fmt.Println("       build kernel")
	return 0, nil
}

func buildHandler(args []string, data interface{}) (int, error) {
	opt := data.(*Options)
	fmt.Printf("invoke 'build' of %s\n", opt.arch)
	return 1, nil
}

func buildKernelHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("invoke 'build kernel', args = %v\n", args)
	return 2, nil
}

func main() {
	opt := &Options{arch: "arm"}

	cmdmux.HandleFunc("/", rootHandler)
	cmdmux.HandleFunc("/build", buildHandler)
	cmdmux.HandleFunc("/build/kernel", buildKernelHandler)

	cmdmux.Execute(opt)
}
```

The package uses `HandleFunc()` to add handler for specific sub-command, like
the package `http`. The sub-command is represented by a command-path, like 
"build", "build/kernel".

After adding handlers, invoke `Execute()` to parse the command line, route to
the correct handler and execute it.

The only one parameter of `Execute()` is passed to the parameter `data` of
the handler function. The parameter `args` of handler function is the rest part
of command line stripped off the command-path part.

The results of this example is like:

```
$ test build
invoke 'build' of arm

$ test build kernel optoins one
invoke 'build kernel', args = [options one]

$ test cmd
Usage: build
       build kernel

$ test
Usage: build
       build kernel
```

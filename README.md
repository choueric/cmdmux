# cmdmux

Package cmdmux implements a command parser and router for terminal program.


[![Build Status](https://travis-ci.org/choueric/cmdmux.svg?branch=master)](https://travis-ci.org/choueric/cmdmux)
[![GoDoc](https://godoc.org/github.com/choueric/cmdmux?status.svg)](https://godoc.org/github.com/choueric/cmdmux)


# Overview

In general, there are two styles a terminal program to interact with users.

1. Use -o, -p to specify parameters. Most programmes are in this way.
2. Use sub-commands, like git which uses only one level sub-command.

The first way can be implemented by the `flag` package of Golang, and the
second this package.

The package can:

1. Build a terminal program with various commands easily.
2. Generate a shell completion file (Now only for bash) !
3. Output commands tree.

# Usage

## Build command line

A simple example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/choueric/cmdmux"
)

type Options struct {
	arch string
}

func rootHandler(args []string, data interface{}) (int, error) {
	fmt.Println("Usage:")
	cmdmux.PrintTree(os.Stderr)
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
	cmdmux.HandleFunc("/build/kernel/image", buildKernelHandler)
	cmdmux.HandleFunc("/build/uboot", buildKernelHandler)

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

The results of this example are like:

```
$ test build
invoke 'build' of arm

$ test build kernel optoins one
invoke 'build kernel', args = [options one]

$ test cmd
Usage:
/*
└── build*
    ├── kernel*
    │   └── image*
    └── uboot*
```

## Print Commands Tree

In the above example, the root handler `rootHandler()` invoke `PrintTree()` to
output the commands tree. If one node has handler, it is appended a '*' symbol
after its command name.

## Generate Shell Completion File

After building various commands with `HandleFunc()`, it's time to get a shell
(bash actually) completion file which helps users of your program input
commands easily on terminal.

Below is a example to get such file for `test` program:

```go
// some HandleFunc codes ...

file, err := os.Create("test-completion")
if err != nil {
	fmt.Println(err)
	return
}
defer file.Close()

if err = cmdmux.GenerateCompletion("test", file); err != nil {
	fmt.Println(err)
}
```

Then apply this file:

```sh
$ source ./test-completion
```

# License

See the LICENSE file.

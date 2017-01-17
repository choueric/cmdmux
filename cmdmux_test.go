/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package cmdmux_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/choueric/cmdmux"
)

const (
	BUILD_RET        = 1
	BUILD_KERNEL_RET = 2
	ROOT_RET         = 3
)

func build(args []string, data interface{}) (int, error) {
	return BUILD_RET, nil
}

func build_kernel(args []string, data interface{}) (int, error) {
	return BUILD_KERNEL_RET, nil
}

func root_handler(args []string, data interface{}) (int, error) {
	return ROOT_RET, nil
}

func flagHandler(args []string, data interface{}) (int, error) {
	str := data.(*string)

	if *str == "test" {
		return 0, nil
	} else {
		return 0, errors.New("flag parsed wrongly")
	}
}

// use the built-in variable
func Test_HandleFunc(t *testing.T) {
	cmdmux.HandleFunc("/build", nil)
	cmdmux.HandleFunc("/build/kernel", nil)
	cmdmux.HandleFunc("/build/dtb", nil)
	cmdmux.HandleFunc("/config/def", nil)
	cmdmux.HandleFunc("/config/menu", nil)
	cmdmux.HandleFunc("/install", nil)

	fmt.Println(cmdmux.String())
}

func Test_Invalid(t *testing.T) {
	cmdMux := cmdmux.New()
	err := cmdMux.HandleFunc("build", nil)
	if err == nil {
		t.Error("should be erroneous")
	}
}

func Test_Execute_normal(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest", "build", "kernel"}

	ret, err := cmdMux.Execute(nil)
	if err != nil {
		t.Error(err)
	}

	if ret != BUILD_KERNEL_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
}

func Test_Execute_opts(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest", "build", "kernel"}

	ret, err := cmdMux.Execute(nil)
	if err != nil {
		t.Error(err)
	}

	if ret != BUILD_KERNEL_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
}

func Test_Execute_noMidNode(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest", "build"}

	_, err := cmdMux.Execute(nil)
	if err == nil {
		t.Error("should be error")
	}
}

func Test_Execute_midNode(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build", build)
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest", "build"}

	ret, err := cmdMux.Execute(nil)
	if err != nil {
		t.Error(err)
	}

	if ret != BUILD_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
}

func Test_Execute_root(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/", root_handler)
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest"}

	ret, err := cmdMux.Execute(nil)
	if err != nil {
		t.Error(err)
	}

	if ret != ROOT_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
}

func Test_Completion(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build", nil)
	cmdMux.HandleFunc("/build/kernel", nil)
	cmdMux.HandleFunc("/build/dtb", nil)
	cmdMux.HandleFunc("/config/def", nil)
	cmdMux.HandleFunc("/config/menu", nil)
	cmdMux.HandleFunc("/install", nil)

	os.Args = []string{"gotest", "UUDDLRLRBABA", "completion"}

	file, err := os.Create("gotest-completion")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	if err = cmdMux.OutputCompletion("test", file); err != nil {
		t.Error(err)
	}
}

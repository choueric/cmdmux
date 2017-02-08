package cmdmux

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
	cmdmux.HandleFunc("/", root_handler)
	cmdmux.HandleFunc("/config/def", nil)
	cmdmux.HandleFunc("/build/kernel/image", nil)
	cmdmux.HandleFunc("/build/dtb", nil)
	cmdmux.HandleFunc("/build", build)
	cmdmux.HandleFunc("/build/kernel", build_kernel)
	cmdmux.HandleFunc("/config/menu", nil)
	cmdmux.HandleFunc("/install", nil)

	fmt.Printf("output of PrintTree:\n")
	cmdmux.PrintTree(os.Stdout)

	os.Args = []string{"gotest", "build", "kernel", "-p", "test"}
	ret, err := cmdmux.Execute(nil)
	if err != nil {
		t.Error(err)
		return
	}

	if ret != BUILD_KERNEL_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}

	os.Args = []string{"gotest", "build"}

	ret, err = cmdmux.Execute(nil)
	if err != nil {
		t.Error(err)
		return
	}

	if ret != BUILD_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
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
		return
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
		return
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
		return
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
		return
	}

	if ret != ROOT_RET {
		t.Errorf("return value wrong: %d\n", ret)
	}
}

func Test_Completion(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build", nil)
	cmdMux.HandleFunc("/build/kernel", nil)
	cmdMux.HandleFunc("/build/kernel/image", nil)
	cmdMux.HandleFunc("/build/dtb", nil)
	cmdMux.HandleFunc("/config/def", nil)
	cmdMux.HandleFunc("/config/menu", nil)
	cmdMux.HandleFunc("/install", nil)

	file, err := os.Create("gotest-completion")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	if err = cmdMux.GenerateCompletion("test", file); err != nil {
		t.Error(err)
	}
}

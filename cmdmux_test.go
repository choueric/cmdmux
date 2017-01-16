package cmdmux_test

import (
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
	os.Args = []string{"gotest", "build", "kernel", "-p", "tk1"}

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
	os.Args = []string{"gotest", "build", "-p", "tk1"}

	_, err := cmdMux.Execute(nil)
	if err == nil {
		t.Error("should be error")
	}
}

func Test_Execute_midNode(t *testing.T) {
	cmdMux := cmdmux.New()
	cmdMux.HandleFunc("/build", build)
	cmdMux.HandleFunc("/build/kernel", build_kernel)
	os.Args = []string{"gotest", "build", "-p", "tk1"}

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

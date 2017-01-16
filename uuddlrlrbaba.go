package cmdmux

import "fmt"

func getCompletionHandler(args []string, data interface{}) (int, error) {
	fmt.Println("getCompletionHandler")
	return 0, nil
}

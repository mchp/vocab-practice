package translate

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no argument supplied")
		return
	}
	if len(args) > 1 {
		fmt.Printf("superfluous arguments %v\n", args[1:])
	}
	words, err := Lookup(args[0])
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	var result []string
	for _, w := range words {
		result = append(result, w.translation)
	}
	fmt.Printf("result: %v\n", result)
}

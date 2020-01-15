package main

import (
	"fmt"
	"os"
	"vocabpractice/data"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no argument supplied")
		return
	}
	local := false
	if args[0] == "local" {
		if len(args) == 1 {
			args = append(args, "none")
		}
		local = true
		args = args[1:]
	}
	db, err := data.InitDynamoDB(local)

	if err != nil {
		fmt.Print("init db:", err)
		return
	}
	switch args[0] {
	case "fetch":
		fmt.Println(db.FetchNext())
	case "query":
		fmt.Println(db.QueryWord(args[1]))
	case "pass":
		fmt.Println(db.Pass(args[1], args[2]))
	case "input":
		fmt.Println(db.Input(args[1], args[2]))
	default:
		fmt.Println("invalid input, try fetch, query, pass or input")
	}
}

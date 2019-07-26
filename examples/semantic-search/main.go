package main

import (
	"bufio"
	"fmt"
	"os"
)

const ExitText = "exit"

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Incorrect args\nUsage: search MODELPATH CSVPATH")
		os.Exit(1)
	}
	modelPath := args[1]
	csvPath := args[2]
	e, err := newEngine(modelPath, csvPath)
	if err != nil {
		panic(err)
	}
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Printf("Engine Initialized\n\n")
	fmt.Println("Enter Query or \"exit\":")
	for stdin.Scan() {
		q := stdin.Text()
		switch q {
		case ExitText:
			fmt.Println("Exiting...")
			return
		case "":
		default:
			res, err := e.search(q)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Suggestion: %+v\n\n", res)
		}
		fmt.Println("Enter Query or \"exit\":")
	}
	if stdin.Err() != nil {
		panic(err)
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
)

const ExitText = "exit"

// Using a convention for this project that _* is a cmdline arg
var (
	_batch       int
	_seqlen      int
	_delim       string
	_d           rune
	_workerCount int
	_modelPath   string
	_csvPath     string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of search:\nArgs: MODELPATH CSVPATH\n")
		flag.PrintDefaults()
	}
	flag.IntVar(&_batch, "b", 32, "Size of batch to encode")
	flag.IntVar(&_seqlen, "seqlen", 16, "Max sequence length")
	flag.StringVar(&_delim, "d", ",", `CSV delimiter char, ex -d=\t`)
	flag.IntVar(&_workerCount, "w", runtime.NumCPU(), "Number of workers for prediction")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		exit("Error: Incorrect args, requires exactly 2 - ", args)
	} else if _delim == "\t" || _delim == "t" {
		fmt.Fprintf(os.Stderr, "Warning: Setting delimiter to tab char\n")
		_d = '\t'
	} else if len(_delim) > 1 {
		exit(`Error: Delimiter set to char t, did you mean -d='\t'`)
	} else {
		_d = rune(_delim[0])
	}
	_modelPath = args[0]
	_csvPath = args[1]
}

func main() {
	e, err := newEngine(_modelPath, int32(_seqlen))
	if err != nil {
		exit("Error:", err)
	}
	if err = e.loadCSV(_csvPath, _d); err != nil {
		exit("Error:", err)
	}
	fmt.Println(e.recs)
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
				exit("Error:", err)
			}
			fmt.Printf("Suggestion: %+v\n\n", res)
		}
		fmt.Println("Enter Query or \"exit\":")
	}
	if stdin.Err() != nil {
		exit("Error:", err)
	}
}

func exit(msgs ...interface{}) {
	flag.Usage()
	fmt.Fprintln(os.Stderr, msgs...)
	os.Exit(1)
}

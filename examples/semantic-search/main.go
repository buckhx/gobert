package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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
	stdin := bufio.NewScanner(os.Stdin)
	log.Printf("Engine Initialized\n\n")
	fmt.Printf("Enter Query or \"exit\":\n\n")
	for stdin.Scan() {
		q := stdin.Text()
		switch q {
		case ExitText:
			fmt.Println("Exiting...")
			return
		case "":
		default:
			res, score, err := e.search(q)
			if err != nil {
				exit("Error:", err)
			}
			fmt.Printf("-> %s\n", res[TextHeader])
			fmt.Printf("\tSimilarity Score (%.2f)\n", score)
			if score < 0.9 {
				fmt.Println("\tNot so sure about that, might need to look somewhere else...")
			} else {
				fmt.Println("\tLGTM")
			}
		}
		fmt.Printf("\nEnter Query or \"exit\":\n\n")
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

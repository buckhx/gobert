package main

import (
	"fmt"
	//	"sync"
	//	"time"

	"github.com/buckhx/gobert/model"
)

func main() {
	path := "./python/output"
	vocabPath := "./bert-models/uncased_L-12_H-768_A-12/vocab.txt"
	m, err := model.NewBert(path, vocabPath, 64)
	if err != nil {
		panic(err)
	}
	infs, err := m.Infer("the dog is hairy.")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inference Count: %d\n", len(infs))
}

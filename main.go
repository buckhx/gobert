package main

import (
	"fmt"
	//	"sync"
	//	"time"

	"github.com/buckhx/gobert/model"
)

func main() {
	path := "/tmp/model/bert/export/output"
	m, err := model.NewBert(path)
	if err != nil {
		panic(err)
	}
	res, err := m.PredictValues("the dog is hairy.")
	if err != nil {
		panic(err)
	}
	vals := res[0].Value().([][][]float32)
	fmt.Printf("Inference Count: %d - %v\n", len(res), vals[0][0][0:16])
}

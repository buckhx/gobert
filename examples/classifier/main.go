package main

import (
	"fmt"
	"os"

	"github.com/buckhx/gobert/model"
)

/*
1. Download base model
2. Fine tune w/ run_classifier
3. export_classifier $MODEL_DIR $EXPORT_DIR 2
4. GOBERT_BASE_DIR=$EXPORT_DIR go run main.go

*/
func main() {
	path := os.Getenv("GOBERT_BASE_DIR")
	m, err := model.NewBertClassifier(path, path+"/vocab.txt")
	if err != nil {
		panic(err)
	}
	texts := []string{
		"the dog that I own is hairy ||| my dog is hairy",
		"there are a lot of bears ||| watch out for bears!",
		"fireworks are for the 4th of july ||| independence day is reason fireworks were created",
		"fireworks are for the 4th of july ||| the fourth of july is reason fireworks were created",
	}
	res, err := m.PredictValues(texts...)
	if err != nil {
		panic(err)
	}
	fmt.Println(res[0].Value().([][]float32))
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/buckhx/gobert/model"
)

/*
1. Download base model
2. Fine tune w/ run_classifier
3. export_classifier $MODEL_DIR $EXPORT_DIR 2
4. MODEL_PATH=$EXPORT_DIR go run main.go

*/
func main() {
	path := os.Getenv("MODEL_PATH")
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
	probs := res[0].Value().([][]float32)
	for i, text := range texts {
		pairs := strings.Split(text, " ||| ")
		same := probs[i][1]
		msg := "Unsure"
		if same > 0.8 {
			msg = "Same"
		} else if same < 0.2 {
			msg = "Different"
		}
		fmt.Println("Meaning:", msg)
		fmt.Printf("\t%q\n\t%q\n\t%v\n", pairs[0], pairs[1], probs[i])
	}
}

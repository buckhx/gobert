package main

import (
	"fmt"
	"os"

	"github.com/buckhx/gobert/model"
)

func main() {
	path := os.Getenv("GOBERT_BASE_DIR")
	m, err := model.NewEmbeddings(path)
	if err != nil {
		panic(err)
	}
	texts := []string{"the dog is hairy.", "watch out for bears!"}
	res, err := m.PredictValues(texts...)
	if err != nil {
		panic(err)
	}
	vals := res[0].Value().([][][]float32)
	fs := m.Features(texts...)
	fmt.Println("\nWord Embeddings")
	for i, s := range vals {
		fmt.Println(texts[i], fs[i].ID)
		for j, toks := range s {
			if toks[0] == 0 {
				break // Hack to not print whole vector
			}
			//fmt.Printf("\t%s\t[%d]float32\t%v...\n", fs[i].Tokens[j], len(toks), toks[:8])
			fmt.Printf("\t%s\t%v...\n", fs[i].Tokens[j], toks[:8])
		}
	}
}

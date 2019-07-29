package main

import (
	"fmt"
	"os"

	"github.com/buckhx/gobert/model"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func main() {
	path := os.Getenv("MODEL_PATH")
	m, err := model.NewEmbeddings(path)
	if err != nil {
		panic(err)
	}
	texts := []string{
		"yeah he's leaving on the midnight train to georgia",
		"the dog is hairy.",
		"watch out for bears!",
		"the cat is red",
		"in 1492 columbus sailed the ocean blue",
		"champa is coloring some stuff",
		"patrick took a train across europe",
	}
	for n := 0; n < 5; n++ {
		fmt.Println("Prediting Vals...")
		res, err := m.PredictValues(texts...)
		fmt.Println("Done Predicting.")
		if err != nil {
			panic(err)
		}
		vals := res[0].Value().([][][]float32)
		fmt.Println("Pooling...")
		embs := make([]mat.Vector, len(vals))
		for s, sent := range vals {
			vec := pool(sent)
			embs[s] = vec
		}
		fmt.Println("Done Pooling.")
		for i := 1; i < len(vals); i++ {
			fmt.Printf("%q, %q -> %v\n", texts[0], texts[i], cosSim(embs[0], embs[i]))
		}
	}
}

func pool(toks [][]float32) mat.Vector {
	c := len(toks[0])
	vec := mat.NewVecDense(c, nil)
	x := make([]float64, c)
	for i := range x {
		for j, tok := range toks {
			x[j] = float64(tok[i])
		}
		vec.SetVec(i, stat.Mean(x, nil))
	}
	return vec
}

func cosSim(x, y mat.Vector) float64 {
	return (mat.Dot(x, y)) / (mat.Norm(x, 2) * mat.Norm(y, 2))

}

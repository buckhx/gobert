package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/buckhx/gobert/model"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type engine struct {
	mod  model.Bert
	recs []map[string]string
	vecs []mat.Vector
}

func newEngine(modelPath, csvPath string, d rune) (engine, error) {
	recs, err := readCSV(csvPath, d)
	if err != nil {
		return engine{}, err
	}
	c := 0
	texts := make([]string, len(recs))
	for i, rec := range recs {
		texts[i] = rec[TextHeader]
		c += len(strings.Split(rec[TextHeader], " "))
	}
	fmt.Println("Average tokens per text", (c/len(texts))+2)         // Account for CLS/SEP
	mod, err := model.NewEmbeddings(modelPath, model.WithSeqLen(16)) // TODO config, avg 11 in quora
	if err != nil {
		return engine{}, err
	}
	var vecs []mat.Vector
	bsize := 16 // TOD Obetter batching
	for to := bsize; to < len(texts)-bsize; to += bsize {
		from := to - bsize
		log.Printf("Predicting Batch Size %d [%d,%d)", bsize, from, to)
		batch := texts[from:to]
		res, err := mod.PredictValues(batch...)
		if err != nil {
			return engine{}, err
		}
		vals := res[0].Value().([][][]float32)
		for _, sent := range vals {
			vecs = append(vecs, MeanPool(sent))
		}
	}
	//TODO submit final batch
	return engine{
		mod:  mod,
		recs: recs,
		vecs: vecs,
	}, nil
}

func (e engine) search(text string) (map[string]string, error) {
	res, err := e.mod.PredictValues(text)
	if err != nil {
		return nil, err
	}
	qvec := MeanPool(res[0].Value().([][][]float32)[0])
	idx := -1
	score := 0.0
	for i, vec := range e.vecs {
		sim := CosSim(qvec, vec)
		if sim > score {
			idx = i
			score = sim
		}
	}
	return e.recs[idx], nil
}

func MeanPool(toks [][]float32) mat.Vector {
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

func CosSim(x, y mat.Vector) float64 {
	return (mat.Dot(x, y)) / (mat.Norm(x, 2) * mat.Norm(y, 2))

}

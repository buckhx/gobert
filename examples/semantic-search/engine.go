package main

import (
	"fmt"
	"github.com/buckhx/gobert/model"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"strings"
)

type engine struct {
	mod  model.Bert
	recs []map[string]string
	vecs []mat.Vector
}

func newEngine(modelPath, csvPath string) (engine, error) {
	recs, err := readCSV(csvPath)
	if err != nil {
		return engine{}, err
	}
	c := 0
	texts := make([]string, len(recs))
	for i, rec := range recs {
		texts[i] = rec[TextHeader]
		c += len(strings.Split(rec[TextHeader], " "))
	}
	mod, err := model.NewEmbeddings(modelPath, model.WithSeqLen(16)) // TODO config, avg 11 in quora
	if err != nil {
		return engine{}, err
	}
	var vecs []mat.Vector
	bsize := 32 // TOD Obetter batching
	for to := bsize; to < len(texts)-bsize; to += bsize {
		fmt.Println("Items", to)
		from := to - bsize
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

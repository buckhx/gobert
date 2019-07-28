package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/buckhx/gobert/model"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

// engine is a simple semantic search engine for demonstrating using a BERT model
// It is not exported because it is not meant to be used outside of demonstration purposes
type engine struct {
	mod  model.Bert
	recs []map[string]string
	vecs []mat.Vector
}

func newEngine(modelPath string, seqlen int32) (*engine, error) {
	mod, err := model.NewEmbeddings(modelPath,
		model.WithSeqLen(seqlen),
	)
	if err != nil {
		return nil, err
	}
	return &engine{
		mod: mod,
	}, nil
}

func (e *engine) loadCSV(csvPath string, d rune) error {
	recs, err := readCSV(csvPath, d)
	if err != nil {
		return err
	}
	tc := 0
	texts := make([]string, len(recs))
	for i, rec := range recs {
		texts[i] = rec[TextHeader]
		tc += len(strings.Split(rec[TextHeader], " "))
	}
	fmt.Println("Average Token Per Text Estimate:", tc/len(texts))
	bsize := _batch                 // TODO batch from flag
	type rng struct{ from, to int } // [from, to)
	ranges := make(chan rng)
	go func() {
		var from, to int
		for from = 0; to < len(texts)-bsize; from = to {
			to += bsize
			ranges <- rng{from: from, to: to}
		}
		ranges <- rng{from: to, to: len(texts)} // final batch
		close(ranges)
	}()
	vecs := make([]mat.Vector, len(texts))
	var workers errgroup.Group
	for i := 0; i < _workerCount; i++ {
		w := i
		workers.Go(func() error {
			for b := range ranges {
				log.Printf("Worker %d - Predicting Batch Size %d [%d,%d)", w, bsize, b.from, b.to)
				batch := texts[b.from:b.to]
				res, err := e.mod.PredictValues(batch...)
				if err != nil {
					return err
				}
				vals := res[0].Value().([][][]float32)
				for i, v := range vals {
					vecs[i+b.from] = meanPool(v)
				}
			}
			return nil
		})
	}
	if err := workers.Wait(); err != nil {
		return err
	}
	e.vecs = append(e.vecs, vecs...)
	e.recs = append(e.recs, recs...)
	return nil
}

func (e *engine) search(text string) (map[string]string, error) {
	res, err := e.mod.PredictValues(text)
	if err != nil {
		return nil, err
	}
	qvec := meanPool(res[0].Value().([][][]float32)[0])
	idx := -1
	score := 0.0
	for i, vec := range e.vecs {
		sim := cosSim(qvec, vec)
		if sim > score {
			idx = i
			score = sim
		}
	}
	return e.recs[idx], nil
}

// TODO extract this into a reusable package
func meanPool(toks [][]float32) mat.Vector {
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

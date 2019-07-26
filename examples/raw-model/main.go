package main

import (
	"fmt"
	"os"

	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/vocab"
)

func main() {
	modelPath := os.Getenv("MODEL_PATH")
	vocabPath := "" + "/vocab.txt"
	voc, err := vocab.FromFile(vocabPath)
	if err != nil {
		panic(err)
	}
	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: 120}
	fs := ff.Feature("the dog is hairy.")
	m, err := tf.LoadSavedModel(modelPath, []string{"bert-untuned"}, nil)
	if err != nil {
		panic(err)
	}
	tids, err := tf.NewTensor([][]int32{f.TokenIDs})
	if err != nil {
		panic(err)
	}
	mask, err := tf.NewTensor([][]int32{f.Mask})
	if err != nil {
		panic(err)
	}
	sids, err := tf.NewTensor([][]int32{f.TypeIDs})
	if err != nil {
		panic(err)
	}
	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("input_ids").Output(0):      tids,
			m.Graph.Operation("input_mask").Output(0):     mask,
			m.Graph.Operation("input_type_ids").Output(0): sids,
		},
		[]tf.Output{
			m.Graph.Operation("embedding").Output(0),
		},
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("DataType", vec[0].DataType())
	fmt.Println("Shape", vec[0].Shape())
	fmt.Println("Value", vec[0].Value().([][][]float32))
}

package main

import (
	"fmt"

	"github.com/buckhx/gobert/model"
	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	path := "./python/output"
	tags := []string{"bert-uncased"}
	vocabPath := "./bert-models/uncased_L-12_H-768_A-12/vocab.txt"
	voc, err := vocab.FromFile(vocabPath)
	if err != nil {
		panic(err)
	}
	tkz := tokenize.NewTokenizer(voc)
	text := "the dog is hairy."
	f := model.SequenceFeature(tkz, 128, text)
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
	m, err := tf.LoadSavedModel(path, tags, nil)
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
			m.Graph.Operation("outputs").Output(0),
		},
		nil,
	)
	if err != nil {
		panic(err)
	}
	vec := res[0]
	fmt.Println("DataType", vec.DataType())
	fmt.Println("Shape", vec.Shape())
	fmt.Println("Value", vec.Value().([][][]float32)[0][0])
}

func debugModel(m *tf.SavedModel) {
	fmt.Printf("%+v\n", m)
	fmt.Println("Session")
	fmt.Println("\tDevice")
	devs, err := m.Session.ListDevices()
	if err != nil {
		fmt.Println(err)
	}
	for _, dev := range devs {
		fmt.Printf("\t\t%+v\n", dev)
	}
	fmt.Println("Graph")
	fmt.Println("\tOperation")
	for _, op := range m.Graph.Operations() {
		fmt.Printf("\t\t%+v\n", op.Name())
	}
}

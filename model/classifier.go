package model

import (
	"fmt"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func NewBertClassifier(path string, vocabPath string, opts ...BertOption) (Bert, error) {
	m, err := tf.LoadSavedModel(path, []string{"bert-tuned"}, nil)
	if err != nil {
		return Bert{}, err
	}
	return NewBert(m, vocabPath, append(opts,
		WithSeqLen(64),
		WithModelFunc(func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			fmt.Println(m.Graph.Operation("loss/Softmax").Output(0).Shape())
			return []tf.Output{
				m.Graph.Operation("loss/Softmax").Output(0),
			}, nil
		}),
	)...)
}

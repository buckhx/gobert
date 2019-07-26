package model

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

const (
	ClassifierOutputOp      = "probabilities"
	ClassifierModelTag      = "bert-tuned"
	DefaultClassifierSeqLen = 64
)

func NewBertClassifier(path string, vocabPath string, opts ...BertOption) (Bert, error) {
	m, err := tf.LoadSavedModel(path, []string{ClassifierModelTag}, nil)
	if err != nil {
		return Bert{}, err
	}
	return NewBert(m, vocabPath, append(opts,
		WithSeqLen(DefaultClassifierSeqLen),
		WithModelFunc(func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			return []tf.Output{
				m.Graph.Operation(ClassifierOutputOp).Output(0),
			}, nil
		}),
	)...)
}

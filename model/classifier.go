package model

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// DefaultOverrides
const (
	ClassifierOutputOp = "probabilities"
	ClassifierModelTag = "bert-tuned"
	ClassifierSeqLen   = 64
)

// NewBertClassifier returns a model configured for classification after being fine-tuned with run_classification.py
func NewBertClassifier(path string, vocabPath string, opts ...BertOption) (Bert, error) {
	m, err := tf.LoadSavedModel(path, []string{ClassifierModelTag}, nil)
	if err != nil {
		return Bert{}, err
	}
	return NewBert(m, vocabPath, append(opts,
		WithSeqLen(ClassifierSeqLen),
		WithModelFunc(func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			return []tf.Output{
				m.Graph.Operation(ClassifierOutputOp).Output(0),
			}, nil
		}),
	)...)
}

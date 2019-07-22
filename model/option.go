package model

import (
	"github.com/buckhx/gobert/model/estimator"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type BertOption func(b Bert) Bert

func WithFeatureFactory(ff *FeatureFactory) BertOption {
	return func(b Bert) Bert {
		b.factory = ff
		return b
	}
}

func WithModelFunc(fn estimator.ModelFunc) BertOption {
	return func(b Bert) Bert {
		b.modelFunc = fn
		return b
	}
}

func WithInputFunc(fn TensorInputFunc) BertOption {
	return func(b Bert) Bert {
		b.inputFunc = fn
		return b
	}
}

func WithSavedModel(m *tf.SavedModel) BertOption {
	return func(b Bert) Bert {
		b.m = m
		return b
	}
}

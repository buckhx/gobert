package model

import (
	"github.com/buckhx/gobert/model/estimator"
	"github.com/buckhx/gobert/tokenize"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type BertOption func(b Bert) Bert

func WithTokenizer(tkz tokenize.VocabTokenizer) BertOption {
	return func(b Bert) Bert {
		b.factory.Tokenizer = tkz
		return b
	}
}

func WithSeqLen(l int32) BertOption {
	return func(b Bert) Bert {
		b.factory.SeqLen = l
		return b
	}
}

func WithFeatureFactory(ff *tokenize.FeatureFactory) BertOption {
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

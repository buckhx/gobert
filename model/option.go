package model

import (
	"github.com/buckhx/gobert/model/estimator"
	"github.com/buckhx/gobert/tokenize"
)

/*
type Bert struct {
	m         *tf.SavedModel
	tokenizer tokenize.VocabTokenizer
	seqLen    int32
	modelFn   estimator.ModelFunc
	inputFn   estimator.InputFunc
}
*/

type BertOption func(b Bert) Bert

func WithSeqLen(l int32) BertOption {
	return func(b Bert) Bert {
		b.seqLen = l
		return b
	}
}

func WithModelFunc(fn estimator.ModelFunc) BertOption {
	return func(b Bert) Bert {
		b.modelFn = fn
		return b
	}
}

func WithInputFunc(fn estimator.InputFunc) BertOption {
	return func(b Bert) Bert {
		b.inputFn = fn
		return b
	}
}

func WithTokenizer(tk tokenize.VocabTokenizer) BertOption {
	return func(b Bert) Bert {
		b.tokenizers = tk
		return b
	}
}

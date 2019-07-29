package model

import (
	"github.com/buckhx/gobert/model/estimator"
	"github.com/buckhx/gobert/tokenize"
)

// BertOption configures a BERT model
type BertOption func(b Bert) Bert

// WithTokenizer applies the given tokenizer to the model
func WithTokenizer(tkz tokenize.VocabTokenizer) BertOption {
	return func(b Bert) Bert {
		b.factory.Tokenizer = tkz
		return b
	}
}

// WithSeqLen applies the seqlen, should match max_seq_len from trained model
func WithSeqLen(l int32) BertOption {
	return func(b Bert) Bert {
		b.factory.SeqLen = l
		return b
	}
}

// WithFeatureFactory replaces the default feature factory
func WithFeatureFactory(ff *tokenize.FeatureFactory) BertOption {
	return func(b Bert) Bert {
		b.factory = ff
		return b
	}
}

// WithModelFunc applies the given model func, used when outputs do not match the default
func WithModelFunc(fn estimator.ModelFunc) BertOption {
	return func(b Bert) Bert {
		b.modelFunc = fn
		return b
	}
}

// WithInputFunc updates the input func, used if input tensors vary from defaults
func WithInputFunc(fn TensorInputFunc) BertOption {
	return func(b Bert) Bert {
		b.inputFunc = fn
		return b
	}
}

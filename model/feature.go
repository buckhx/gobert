package model

import (
	"strings"
	"sync"

	"github.com/buckhx/gobert/tokenize"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Static tokens
const (
	ClassToken        = "[CLS]"
	SeparatorToken    = "[SEP]"
	SequenceSeparator = " ||| "
)

// Feature is an input feature for a BERT model.
// Maps to extract_features.InputFeature in ref-impl
type Feature struct {
	ID       int32
	Text     string
	Tokens   []string
	TokenIDs []int32
	Mask     []int32 // short?
	TypeIDs  []int32 // seqeuence ids, short?
	//	LabelIDs []int32
}

// Size will return the number of tokens in the feature by counting the mask bits
func (f Feature) Size() int {
	var l int
	for _, v := range f.Mask {
		if v > 0 {
			l += 1
		}
	}
	return l
}

func tensors(fs ...Feature) (map[string]*tf.Tensor, error) {
	//	uids := make([]int32, len(fs))
	tids := make([][]int32, len(fs))
	mask := make([][]int32, len(fs))
	sids := make([][]int32, len(fs))
	for i, f := range fs {
		//		uids[i] = f.ID
		tids[i] = f.TokenIDs
		mask[i] = f.Mask
		sids[i] = f.TypeIDs
	}
	/*
		u, err := tf.NewTensor(uids)
		if err != nil {
			return nil, err
		}
	*/
	t, err := tf.NewTensor(tids)
	if err != nil {
		return nil, err
	}
	m, err := tf.NewTensor(mask)
	if err != nil {
		return nil, err
	}
	s, err := tf.NewTensor(sids)
	if err != nil {
		return nil, err
	}
	return map[string]*tf.Tensor{
		//UniqueIDsOp:    u,
		InputIDsOp:     t,
		InputMaskOp:    m,
		InputTypeIDsOp: s,
	}, nil
}

type FeatureFactory struct {
	tokenizer tokenize.VocabTokenizer
	seqLen    int32
	count     int32
	lock      sync.Mutex
}

func (ff *FeatureFactory) Feature(text string) Feature {
	f := sequenceFeature(ff.tokenizer, ff.seqLen, text)
	ff.lock.Lock()
	f.ID = ff.count
	ff.count += 1
	ff.lock.Unlock()
	return f
}

func (ff *FeatureFactory) Features(texts ...string) []Feature {
	fs := make([]Feature, len(texts))
	for i, text := range texts {
		fs[i] = ff.Feature(text)
	}
	return fs
}

// SequenceFeature will take a sequence string and
// build features for the model from it
func sequenceFeature(tkz tokenize.VocabTokenizer, seqLen int32, text string) Feature {
	f := Feature{
		Text:     text,
		Tokens:   make([]string, seqLen),
		TokenIDs: make([]int32, seqLen),
		Mask:     make([]int32, seqLen),
		TypeIDs:  make([]int32, seqLen),
	}
	parts := strings.Split(text, SequenceSeparator)
	seqs := make([][]string, len(parts))
	for i, part := range parts {
		seqs[i] = tkz.Tokenize(part)
	}
	seqs = truncate(seqs, seqLen-int32(len(seqs))-1) // truncate w/ space for CLS/SEP
	voc := tkz.Vocab()
	var s int
	f.Tokens[s] = ClassToken
	f.TokenIDs[s] = voc.GetID(ClassToken).Int32()
	f.TypeIDs[s] = 0
	f.Mask[s] = 1
	s++
	for sid, seq := range seqs {
		for _, tok := range seq {
			f.Tokens[s] = tok
			f.TokenIDs[s] = voc.GetID(tok).Int32()
			f.TypeIDs[s] = int32(sid)
			f.Mask[s] = 1
			s++
		}
		f.Tokens[s] = SeparatorToken
		f.TokenIDs[s] = voc.GetID(SeparatorToken).Int32()
		f.TypeIDs[s] = int32(sid)
		f.Mask[s] = 1
		s++
	}
	return f
}

// truncate uses heuristic of trimming seq with longest len until seqlen satisfied
func truncate(seqs [][]string, maxlen int32) [][]string {
	// TODO test
	// NOTE: this impl could be a bottleneck
	var seqlen int32
	for i := range seqs {
		seqlen += int32(len(seqs[i]))
	}
	for slen := seqlen; slen > maxlen; slen-- {
		// Sort to get longest first
		var mi, mv int
		for i := len(seqs) - 1; i >= 0; i-- {
			seq := seqs[i] // iterate in reverse to select lower indexes
			if len(seq) > mv {
				mi = i
				mv = len(seq)
			}
		}
		if mv <= 0 { // can't trim anymore
			return seqs
		}
		rm := seqs[mi]
		rm[len(rm)-1] = "" // Mark for GC, avoid mem leak
		seqs[mi] = rm[:len(rm)-1]
	}
	return seqs
}

package tokenize

import (
	"strings"
	"sync"
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
}

// Count will return the number of tokens in the feature by counting the mask bits
func (f Feature) Count() int {
	var l int
	for _, v := range f.Mask {
		if v > 0 {
			l++
		}
	}
	return l
}

// FeatureFactory will create features with the supplied tokenizer and sequence length
type FeatureFactory struct {
	Tokenizer VocabTokenizer
	SeqLen    int32
	lock      sync.Mutex
	count     int32
}

// Feature will create a single feature from the factory
// ID creation is thread safe and incremental
func (ff *FeatureFactory) Feature(text string) Feature {
	f := sequenceFeature(ff.Tokenizer, ff.SeqLen, text)
	ff.lock.Lock()
	f.ID = ff.count
	ff.count++
	ff.lock.Unlock()
	return f
}

// Features will create multiple features with incremental IDs
func (ff *FeatureFactory) Features(texts ...string) []Feature {
	fs := make([]Feature, len(texts))
	for i, text := range texts {
		fs[i] = ff.Feature(text)
	}
	return fs
}

// SequenceFeature will take a sequence string and
// build features for the model from it
func sequenceFeature(tkz VocabTokenizer, seqLen int32, text string) Feature {
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

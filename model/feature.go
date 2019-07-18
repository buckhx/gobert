package model

import (
	"strings"

	"github.com/buckhx/gobert/tokenize"
)

const (
	ClassToken        = "[CLS]"
	SeparatorToken    = "[SEP]"
	SequenceSeparator = " ||| "
)

type Feature struct {
	ID       int
	Tokens   []string
	TokenIDs []int32
	Mask     []int // short?
	TypeIDs  []int // seqeuence ids, short?
}

// SequenceFeature will take a sequence string and
// build features for the model from it
func SequenceFeature(tkz tokenize.VocabTokenizer, seqLen int, text string) Feature {
	f := Feature{
		ID:       0, // TODO
		Tokens:   make([]string, seqLen),
		TokenIDs: make([]int32, seqLen),
		Mask:     make([]int, seqLen),
		TypeIDs:  make([]int, seqLen),
	}
	parts := strings.Split(text, SequenceSeparator)
	seqs := make([][]string, len(parts))
	for i, part := range parts {
		seqs[i] = tkz.Tokenize(part)
	}
	seqs = truncate(seqs, seqLen-len(seqs)-1) // truncate w/ space for CLS/SEP
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
			f.TypeIDs[s] = sid
			f.Mask[s] = 1
			// TODO f.TokenIDs
			s++
		}
		f.Tokens[s] = SeparatorToken
		f.TokenIDs[s] = voc.GetID(SeparatorToken).Int32()
		f.TypeIDs[s] = sid
		f.Mask[s] = 1
		s++
	}
	return f
}

// truncate uses heuristic of trimming seq with longest len until seqlen satisfied
func truncate(seqs [][]string, maxlen int) [][]string {
	// TODO test
	// NOTE: this impl could be a bottleneck
	var seqlen int
	for i := range seqs {
		seqlen += len(seqs[i])
	}
	for seqlen = seqlen; seqlen > maxlen; seqlen-- {
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

// package Tokenization supplies tokenizzation operations for BERT.
// NOTE: All defintions are related to BERT and may vary from unicode defintions,
// for example, BERT considers '$' punctuation, but unicode does not.
package tokenization

import "unicode"

// _Bp is the BERT extension of the Punctuation character range
var _Bp = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x002f, 1}, // 33-47
		{0x003a, 0x0040, 1}, // 58-64
		{0x005b, 0x0060, 1}, // 91-96
		{0x007b, 0x007e, 1}, // 123-126
	},
	LatinOffset: 4, // All less than 0x00FF
}

// IsWhitespace checks whether rune c is a BERT whitespace character
func IsWhitespace(c rune) bool {
	switch c {
	case ' ':
		return true
	case '\t':
		return true
	case '\n':
		return true
	case '\r':
		return true
	}
	return unicode.Is(unicode.Zs, c)
}

// IsControl checks wher rune c is a BERT control character
func IsControl(c rune) bool {
	switch c {
	case '\t':
		return false
	case '\n':
		return false
	case '\r':
		return false
	}
	return unicode.In(c, unicode.Cc, unicode.Cf)
}

// IsPunctuation checks wher rune c is a BERT punctuation character
func IsPunctuation(c rune) bool {
	return unicode.In(c, _Bp, unicode.P)
}

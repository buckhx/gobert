// package Tokenization supplies tokenizzation operations for BERT.
// NOTE: All defintions are related to BERT and may vary from unicode defintions,
// for example, BERT considers '$' punctuation, but unicode does not.
package tokenization

import "unicode"

// IsWhitespace checks whether rune c is a whitespace character
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

// IsControl checks wher rune c is a control character
func IsControl(c rune) bool {
	return false

}

func IsPunctuation(c rune) bool {
	return false
}

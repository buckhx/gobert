package tokenization

// ID is used to identify vocab items
type ID int32

/*
func (id ID) String() string {
	return fmt.Sprint(id)
}
*/

// Vocab is a container for tokens
// NOTE: python uses an OrderedDict, unsure of implications
type Vocab struct {
	tokens map[string]ID
}

func NewVocab(tokens []string) Vocab {
	v := make(map[string]ID, len(tokens))
	for i, t := range tokens {
		v[t] = ID(i)
	}
	return Vocab{tokens: v}
}

// Add will add an item to the vocabulary, is not thread-safe
func (v Vocab) Add(token string) {
	v.tokens[token] = ID(v.Size())
}

// Get will return the ID of the token in the vocab. Will be negative if it doesn't exists
func (v Vocab) Get(token string) ID {
	id, ok := v.tokens[token]
	if !ok {
		return ID(-1)
	}
	return ID(id)
}

// SIze returns the size of the vocabulary
func (v Vocab) Size() int {
	return len(v.tokens)
}

// LongestSubstring returns the longest token that is a substring of the token
func (v Vocab) LongestSubstring(token string) string {
	// Greedt, optimize to trie if needed
	for i := len(token); i > 0; i-- {
		sub := token[:i]
		if _, ok := v.tokens[sub]; ok {
			return sub
		}
	}
	return ""
}

func (v Vocab) ConvertItems(items []string) []ID {
	ids := make([]ID, len(items))
	for i, m := range items {
		ids[i] = v.tokens[m]
	}
	return ids
}

func (v Vocab) ConvertTokens(tokens []string) []ID {
	return v.ConvertItems(tokens)
}

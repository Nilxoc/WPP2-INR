package index

import "g1.wpp2.hsnr/inr/boolret/index/spell"

type KGramIndex struct {
	k       int
	entries map[string][]*IndexEntry
}

func InitKGramIndex(k int) *KGramIndex {
	return &KGramIndex{
		k:       k,
		entries: make(map[string][]*IndexEntry),
	}
}

func (idx *KGramIndex) AddKGram(token string, ref *IndexEntry) {
	grams := spell.ExtractKGrams(token, idx.k)
	for _, g := range grams {
		idx.addGram(g, ref)
	}
}

func (idx *KGramIndex) addGram(gram string, ref *IndexEntry) {
	if refList, ok := idx.entries[gram]; ok {
		idx.entries[gram] = append(refList, ref)
	} else {
		idx.entries[gram] = make([]*IndexEntry, 1)
		idx.entries[gram][0] = ref
	}
}

func (idx *KGramIndex) FindTokens(token string) []*IndexEntry {
	needles := spell.ExtractKGrams(token, idx.k)

	res := make([]*IndexEntry, 0)
	for gram, kg := range idx.entries {
		for _, needle := range needles {
			if needle == gram {
				res = append(res, kg...)
			}
		}
	}
	return res
}

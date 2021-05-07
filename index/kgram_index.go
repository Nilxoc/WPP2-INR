package index

import "g1.wpp2.hsnr/inr/boolret/index/spell"

type KGramIndex struct {
	K       int
	Entries map[string][]*IndexEntry
}

type FindTokensRes struct {
	Entry *IndexEntry
	Count int
}

func InitKGramIndex(k int) *KGramIndex {
	return &KGramIndex{
		K:       k,
		Entries: make(map[string][]*IndexEntry),
	}
}

func (idx *KGramIndex) AddKGram(token string, ref *IndexEntry) {
	grams := spell.ExtractKGrams(token, idx.K)
	for _, g := range grams {
		idx.addGram(g, ref)
	}
}

func (idx *KGramIndex) addGram(gram string, ref *IndexEntry) {
	if refList, ok := idx.Entries[gram]; ok {
		idx.Entries[gram] = append(refList, ref)
	} else {
		idx.Entries[gram] = make([]*IndexEntry, 1)
		idx.Entries[gram][0] = ref
	}
}

func (idx *KGramIndex) FindTokens(token string) []*FindTokensRes {
	needles := spell.ExtractKGrams(token, idx.K)

	acc := make(map[string]*FindTokensRes, 0)
	/**
	for gram, kg := range idx.Entries {
		for _, needle := range needles {
			if needle == gram {
				res = append(res, kg...)
			}
		}
	}
	**/

	for _, needle := range needles {
		out, found := idx.Entries[needle]
		if found {
			for _, entry := range out {
				resTerm, resFound := acc[entry.Term]
				if resFound {
					(*resTerm).Count += 1
				} else {
					acc[entry.Term] = &FindTokensRes{entry, 1}
				}
			}
		}
	}

	res := make([]*FindTokensRes, 0)
	for _, val := range acc {
		res = append(res, val)
	}
	return res
}

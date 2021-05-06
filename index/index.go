package index

import (
	"fmt"
	"sort"

	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index/spell"
)

type IndexEntry struct {
	Term string
	Docs PostingList
}

type Index struct {
	Index map[string]*IndexEntry
	K     int
	r     int
	j     float32
	kgram *KGramIndex
}

func NewIndexEmpty(k int, r int, J float32) *Index {
	idx := &Index{
		K: k,
		r: r,
		j: J,
	}
	idx.Index = make(map[string]*IndexEntry)
	idx.kgram = InitKGramIndex(k)
	return idx
}

func NewIndexFromFile(path string, k int, r int, j float32) (*Index, error) {
	idx, err := loadIndex(path)
	idx.r = r
	idx.j = j
	idx.K = k
	idx.kgram = InitKGramIndex(k)
	if err != nil {
		return nil, err
	}

	//Regenerate KGram-Index on Load, because we cannot persist it..
	for k, v := range idx.Index {
		idx.kgram.AddKGram(k, v)
	}
	return idx, nil
}

///INDEX SPECIFIC METHODS
func (i *Index) AddTerm(term string, posting *Posting) {
	if entry, found := i.Index[term]; found {
		entry.Docs = append(entry.Docs, *posting)
		i.kgram.AddKGram(term, entry)
	} else {
		t := IndexEntry{
			Term: term,
			Docs: make(PostingList, 1),
		}
		t.Docs[0] = *posting
		i.Index[term] = &t
		i.kgram.AddKGram(term, &t)
	}
}

func docsSum(inp []*IndexEntry) int {
	sum := 0
	for _, i := range inp {
		sum += len(i.Docs)
	}
	return sum
}

func (i *Index) GetTerm(term string) *IndexEntry {
	if e, f := i.Index[term]; f {
		return e
	}
	return nil
}

func (i *Index) GetTermSuggestions(term string) []*IndexEntry {
	if term == "" {
		return nil
	}

	res := make([]*IndexEntry, 1)
	if e, f := i.Index[term]; f {
		res[0] = e
		docs := docsSum(res)
		if docs < i.r {
			//Not enough documents found
			res = append(res, i.getCorrectedDocs(term, i.r-docs)...)
		}
		return res
	}
	return i.getCorrectedDocs(term, i.r)
}

func dedupe(entries []*IndexEntry) []*IndexEntry {
	res := make([]*IndexEntry, 0, len(entries))
	contains := func(x *IndexEntry) bool {
		for _, e := range res {
			if e.Term == x.Term {
				return true
			}
		}
		return false
	}

	for _, e := range entries {
		if !contains(e) {
			res = append(res, e)
		}
	}
	return res
}

func (i *Index) getCorrectedDocs(term string, altCount int) []*IndexEntry {

	type Candidate struct {
		ldist int
		entry *IndexEntry
	}

	candidates := make([]Candidate, 0)

	posTokens := i.kgram.FindTokens(term)
	posTokens = dedupe(posTokens)

	for _, k := range posTokens {
		if k.Term != term { // SKIP ALREADY SELECED TERM
			if jv := spell.Jaccard(term, k.Term, i.K); jv > i.j {
				candidates = append(candidates, Candidate{
					ldist: spell.LevenshteinDistance(term, k.Term),
					entry: k,
				})
			}
		}
	}

	//Sort by Distance & by string
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ldist == candidates[j].ldist {
			return candidates[i].entry.Term < candidates[j].entry.Term
		}
		return candidates[i].ldist < candidates[j].ldist
	})

	minFunc := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	min := minFunc(len(candidates), altCount)

	res := make([]*IndexEntry, min)
	for i := 0; i < min; i++ {
		fmt.Printf("Did you mean %s ?\n", candidates[i].entry.Term)
		res[i] = candidates[i].entry
	}

	return res
}

func (i *Index) GetTermListCorrected(term string) PostingList {
	tmp := i.GetTermSuggestions(term)
	if len(tmp) == 1 {
		return tmp[0].Docs
	}

	res := make(PostingList, 0)
	for _, e := range tmp {
		res = append(res, e.Docs...)
	}
	return res
}

func (i *Index) SaveIndex(path string) error {
	return file.SaveIndex(i, path)
}

func (i *Index) Len() int {
	return len(i.Index)
}

func (i *IndexEntry) String() string {
	return i.Docs.String()
}

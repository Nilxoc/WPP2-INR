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
	index  map[string]*IndexEntry
	k      int
	r      int
	j      float32
	kgramm *KGramIndex
}

func NewIndexEmpty(k int, r int, J float32) *Index {
	idx := &Index{
		k: k,
		r: r,
		j: J,
	}
	idx.index = make(map[string]*IndexEntry)
	idx.kgramm = InitKGramIndex(k)
	return idx
}

func NewIndexFromFile(path string) (*Index, error) {
	idx, err := loadIndex(path)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

///INDEX SPECIFIC METHODS
func (i *Index) AddTerm(term string, posting *Posting) {
	if entry, found := i.index[term]; found {
		entry.Docs = append(entry.Docs, *posting)
		i.kgramm.AddKGram(term, entry)
	} else {
		t := IndexEntry{
			Term: term,
			Docs: make(PostingList, 1),
		}
		t.Docs[0] = *posting
		i.index[term] = &t
		i.kgramm.AddKGram(term, &t)
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
	if e, f := i.index[term]; f {
		return e
	}
	return nil
}

func (i *Index) GetTermSuggestions(term string) []*IndexEntry {
	if term == "" {
		return nil
	}

	res := make([]*IndexEntry, 1)
	if e, f := i.index[term]; f {
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

func (i *Index) getCorrectedDocs(term string, altCount int) []*IndexEntry {

	type Candidate struct {
		ldist int
		entry *IndexEntry
	}

	candidates := make([]Candidate, 0)

	posTokens := i.kgramm.FindTokens(term)

	for _, k := range posTokens {
		if k.Term != term { // SKIP ALREADY SELECED TERM
			if jv := spell.Jaccard(term, k.Term, i.k); jv > i.j {
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
	return len(i.index)
}

func (i *IndexEntry) String() string {
	return i.Docs.String()
}

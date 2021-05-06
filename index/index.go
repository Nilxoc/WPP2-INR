package index

import (
	"g1.wpp2.hsnr/inr/boolret/file"
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

///INDEX SPECIFIC METHODS / PLACEHOLDER
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

func (i *Index) GetTerm(term string) *IndexEntry {
	if e, f := i.index[term]; f {
		return e
	}
	return nil
}

func (i *Index) getCorrectedTerm(term string) *IndexEntry {
	return nil
}

func (i *Index) SaveIndex(path string) error {
	return file.SaveIndex(i, path)
}

func (i *Index) Len() int {
	return len(i.index)
}

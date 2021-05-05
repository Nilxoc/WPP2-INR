package index

import (
	"g1.wpp2.hsnr/inr/boolret/file"
)

type IndexEntry struct {
	Term string
	Docs PostingList
}

type Index map[string]*IndexEntry

func NewIndexEmpty() *Index {
	idx := make(Index)
	return &idx
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
	if entry, found := (*i)[term]; found {
		entry.Docs = append(entry.Docs, *posting)
	} else {
		t := IndexEntry{
			Term: term,
			Docs: make(PostingList, 1),
		}
		t.Docs[0] = *posting
		(*i)[term] = &t
	}
}

func (i *Index) GetTerm(term string) *IndexEntry {
	if e, f := (*i)[term]; f {
		return e
	}
	return nil
}

func (i *Index) SaveIndex(path string) error {
	return file.SaveIndex(i, path)
}

func (i *Index) Len() int {
	return len(*i)
}

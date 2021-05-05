package index

import (
	"fmt"
)

type IndexEntry struct {
	Term string
	Docs PostingList
}

type Index map[string]IndexEntry

func NewIndexEmpty() *Index {
	//TODO: Create & initialize empty index
	return nil
}

func NewIndexFromFile(path string) (*Index, error) {
	//Load Index Dump from file and initialize index with it (skip tokenizer)
	return nil, fmt.Errorf("Not Implemented!")
}

///INDEX SPECIFIC METHODS / PLACEHOLDER
func (i *Index) AddTerm(term string, posting *PostingListEntry) {
	//TODO
}

func (i *Index) SaveIndex(path string) error {
	return fmt.Errorf("Not Implemented")
}

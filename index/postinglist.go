package index

type PostingListEntry struct {
	DocID int64
	Pos   int64
}

type PostingList []PostingListEntry

func (pl *PostingList) String() string {
	return "TODO"
}

package index

type Posting struct {
	DocID int64
	Pos   []int64
}

type PostingList []Posting

func (pl *PostingList) String() string {
	return "TODO"
}

func (pl *PostingList) Intersect(other *PostingList) PostingList {
	answer := make(PostingList, 0)
	var p1, p2 int
	for p1 < len(*pl) && p2 < len(*other) {
		if (*pl)[p1].DocID == (*other)[p2].DocID {
			answer = append(answer, (*pl)[p1])
			p1 += 1
			p2 += 1
		}

	}
	return answer
}

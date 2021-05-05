package index

import "sort"

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
			newPos := append((*pl)[p1].Pos, (*other)[p2].Pos...)
			sort.Slice(newPos, func(p, q int) bool { return newPos[p] < newPos[q] })
			answer = append(answer, Posting{(*pl)[p1].DocID, newPos})
			p1 += 1
			p2 += 1
		} else if (*pl)[p1].DocID < (*other)[p2].DocID {
			p1 += 1
		} else {
			p2 += 1
		}
	}
	return answer
}

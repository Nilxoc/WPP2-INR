package index

import (
	"math"
	"sort"
)

type Posting struct {
	DocID int64
	Pos   []int64
}

func (p *Posting) Merge(other *Posting) Posting {
	newPos := append(p.Pos, other.Pos...)
	sort.Slice(newPos,
		func(p, q int) bool { return newPos[p] < newPos[q] })
	return Posting{p.DocID, newPos}
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
			answer = append(answer, (*pl)[p1].Merge(&(*other)[p2]))
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

func (pl *PostingList) Union(other *PostingList) PostingList {
	answer := make(PostingList, 0)
	var p1, p2 int
	for p1 < len(*pl) && p2 < len(*other) {
		if (*pl)[p1].DocID == (*other)[p2].DocID {
			answer = append(answer, (*pl)[p1].Merge(&(*other)[p2]))
			p1 += 1
			p2 += 1
		} else if (*pl)[p1].DocID < (*other)[p2].DocID {
			answer = append(answer, (*pl)[p1])
			p1 += 1
		} else {
			answer = append(answer, (*other)[p2])
			p2 += 1
		}
	}

	for p1 < len(*pl) {
		answer = append(answer, (*pl)[p1])
		p1 += 1
	}

	for p2 < len(*other) {
		answer = append(answer, (*other)[p2])
		p1 += 1
	}
	return answer
}

func (pl *PostingList) positionalIntersect(other *PostingList, k int64, kCond func(num1, num2, k int64) bool) PostingList {
	answer := make(PostingList, 0)
	var p1, p2 int
	for p1 < len(*pl) && p2 < len(*other) {
		if (*pl)[p1].DocID == (*other)[p2].DocID {
			l := make([]int64, 0)
			outPos := make([]int64, 0)
			inList := make(map[int64]bool)

			var pos1 = (*pl)[p1].Pos
			var pos2 = (*other)[p2].Pos
			var pp1, pp2 int
			for pp1 < len(pos1) {
				for pp2 < len(pos2) {
					if kCond(pos1[pp1], pos2[pp2], k) {
						l = append(l, pos2[pp2])
					} else if pos2[pp2] > pos1[pp1] {
						break
					}
					pp2 += 1
				}

				for len(l) > 0 && int64(math.Abs(float64(l[0])-float64(pos1[pp1]))) > k {
					if len(l) > 1 {
						l = l[1:]
					} else {
						l = []int64{}
					}
				}

				if len(l) > 0 {
					outPos = append(outPos, pos1[pp1])
				}

				for _, val := range l {
					if !inList[val] {
						outPos = append(outPos, val)
						inList[val] = true
					}
				}
				pp1 += 1
			}
			sort.Slice(outPos, func(p, q int) bool { return outPos[p] < outPos[q] })
			answer = append(answer, Posting{(*pl)[p1].DocID, outPos})
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

func (pl *PostingList) PositionalIntersect(other *PostingList, k int64) PostingList {
	return pl.positionalIntersect(other, k, func(num1, num2, k int64) bool {
		return int64(math.Abs(float64(num1)-float64(num2))) <= k
	})
}

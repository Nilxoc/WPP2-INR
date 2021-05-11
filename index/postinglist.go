package index

import (
	"fmt"
	"math"
	"sort"
)

type Posting struct {
	DocID int64
	Pos   []int64
}

func insert(s []int64, k int, vs ...int64) []int64 {
	//https://github.com/golang/go/wiki/SliceTricks#insert
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]int64, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

func (p *Posting) Merge(other *Posting) Posting {
	newPos := append(p.Pos, other.Pos...)
	sort.Slice(newPos,
		func(p, q int) bool { return newPos[p] < newPos[q] })
	return Posting{p.DocID, newPos}
}

type PostingList []Posting

func (pl *PostingList) Empty() bool {
	return len(*pl) == 0
}

func intersectArrays(a []int64, b []int64) []int64 {
	var i, j int
	res := make([]int64, 0)
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			res = append(res, a[i])
			i += 1
			j += 1
		} else if a[i] < b[j] {
			i += 1
		} else {
			j += 1
		}
	}
	return res
}

func (pl *PostingList) String(idx *Index) string {
	out := ""
	//out += fmt.Sprintln("Found the following documents:")

	for _, p := range *pl {
		out += fmt.Sprintf("%s ", idx.GetDocDisplay(p.DocID))
	}
	return out
}

func (p *Posting) String() string {
	return fmt.Sprintf("%d", p.DocID)
}

func (pl *PostingList) Intersect(other *PostingList) *PostingList {
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
	return &answer
}

func (pl *PostingList) Union(other *PostingList) *PostingList {
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
		p2 += 1
	}
	return &answer
}

func (pl *PostingList) positionalIntersect(other *PostingList, k int64, kCond func(num1, num2, k int64) bool) *PostingList {
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

				for len(l) > 0 && !kCond(pos1[pp1], l[0], k) {
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
	return &answer
}

func (pl *PostingList) Proximity(other *PostingList, k int64) *PostingList {
	return pl.positionalIntersect(other, k, func(num1, num2, k int64) bool {
		return int64(math.Abs(float64(num1)-float64(num2))) <= k
	})
}

func (pl *PostingList) PhraseIntersect(others []*PostingList) *PostingList {
	currPl := *pl
	for _, v := range others {
		currPl = *(currPl.positionalIntersect(v, int64(1), func(num1, num2, k int64) bool { return num2-num1 == k }))
	}

	//Add Positions after first to Position List. This is only important for proximity queries
	//Even then, the middle values could technically be omitted
	for i, posting := range currPl {
		for j, pos := range posting.Pos {
			if (j % 2) != 0 {
				continue
			}

			insertCount := len(others) - 1
			toInsert := make([]int64, insertCount)
			for k := 0; k < len(others)-1; k += 1 {
				toInsert[k] = pos - int64(k+1)
			}

			insertPos := 0
			if j != 0 {
				insertPos = j + ((j-1)/2)*insertCount
			}
			currPl[i].Pos = insert(currPl[i].Pos, insertPos, toInsert...)
		}
	}

	return &currPl
}

func (pl *PostingList) Difference(other *PostingList) *PostingList {
	answer := make(PostingList, 0)
	var p1, p2 int
	for p1 < len(*pl) && p2 < len(*other) {
		if (*pl)[p1].DocID == (*other)[p2].DocID {
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
	return &answer
}

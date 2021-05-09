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

func (pl *PostingList) String() string {
	out := ""
	//out += fmt.Sprintln("Found the following documents:")

	for _, p := range *pl {
		out += fmt.Sprintf("%d ", p.DocID)
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
		p1 += 1
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
	var k int
	for i, v := range others {
		res := currPl.positionalIntersect(v, int64(i+1), func(num1, num2, k int64) bool { return num2-num1 == k })
		for _, v2 := range *res {
			for k < len(currPl) && currPl[k].DocID < v2.DocID {
				currPl = append(currPl[:k], currPl[k+1:]...)
			}

			if k >= len(currPl) {
				k -= 1
				break
			}

			if v2.DocID == currPl[k].DocID {
				currPl[k] = Posting{v2.DocID, intersectArrays(currPl[k].Pos, v2.Pos)}
				k += 1
			}
		}
		currPl = currPl[:k]
		k = 0
	}

	//Add Positions after first to Position List. This is only important for proximity queries
	//Even then, the middle values could technically be omitted
	res := make(PostingList, 0, len(currPl))
	for _, v := range currPl {
		var posList []int64
		for _, pos := range v.Pos {
			posList = append(posList, pos)
			for i := int64(0); i < int64(len(others)); i++ {
				posList = append(posList, pos+int64(1)+i)
			}
		}
		res = append(res, Posting{v.DocID, posList})
	}
	return &res
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

package index

import "testing"

func TestIntersect(t *testing.T) {
	var list1 = make(PostingList, 0)
	var list2 = make(PostingList, 0)

	list1 = append(list1, Posting{1, []int64{2, 5, 8}})
	list1 = append(list1, Posting{5, []int64{6, 9, 20}})
	list1 = append(list1, Posting{8, []int64{3, 8, 89}})

	list2 = append(list2, Posting{1, []int64{1, 6, 7}})
	list2 = append(list2, Posting{4, []int64{8, 23, 91}})
	list2 = append(list2, Posting{8, []int64{25, 39, 90}})

	var res = list1.Intersect(&list2)
	if res[0].DocID != 1 {
		t.Errorf("First DocID Wrong, Expected 1, got %d", res[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{1, 2, 5, 6, 7, 8}

	if len(res[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range res[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("Wrong Positions in first posting, expected [1,2,5,6,7,8],got %v", res[0].Pos)
	}

	if res[1].DocID != 8 {
		t.Errorf("Second DocID Wrong, Expected 8, got %d", res[1].DocID)
	}

	arrayEqual = true
	correct = []int64{3, 8, 25, 39, 89, 90}

	if len(res[1].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range res[1].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("Wrong Positions in first second, expected [1,2,5,6,7,8],got %v", res[1].Pos)
	}
}

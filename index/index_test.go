package index

import (
	"testing"
)

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

	if len(*res) != 2 {
		t.Errorf("wrong posting count, expected 2, got %d", len(*res))
	}

	if (*res)[0].DocID != 1 {
		t.Errorf("First DocID Wrong, Expected 1, got %d", (*res)[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{1, 2, 5, 6, 7, 8}

	if len((*res)[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("Wrong Positions in first posting, expected [1,2,5,6,7,8],got %v", (*res)[0].Pos)
	}

	if (*res)[1].DocID != 8 {
		t.Errorf("Second DocID Wrong, Expected 8, got %d", (*res)[1].DocID)
	}

	arrayEqual = true
	correct = []int64{3, 8, 25, 39, 89, 90}

	if len((*res)[1].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[1].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("Wrong Positions in first second, expected [1,2,5,6,7,8],got %v", (*res)[1].Pos)
	}
}

func TestUnion(t *testing.T) {
	var list1 = make(PostingList, 0)
	var list2 = make(PostingList, 0)

	list1 = append(list1, Posting{1, []int64{2, 5, 8}})
	list1 = append(list1, Posting{5, []int64{6, 9, 20}})
	list1 = append(list1, Posting{8, []int64{3, 8, 89}})

	list2 = append(list2, Posting{1, []int64{1, 6, 7}})
	list2 = append(list2, Posting{4, []int64{8, 23, 91}})
	list2 = append(list2, Posting{8, []int64{25, 39, 90}})

	var res = list1.Union(&list2)

	if len(*res) != 4 {
		t.Errorf("wrong posting count, expected 4, got %d", len(*res))
	}

	if (*res)[0].DocID != 1 {
		t.Errorf("First DocID Wrong, Expected 1, got %d", (*res)[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{1, 2, 5, 6, 7, 8}

	if len((*res)[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in first posting, expected [1,2,5,6,7,8],got %v", (*res)[0].Pos)
	}

	if (*res)[1].DocID != 4 {
		t.Errorf("second docid Wrong, expected 4, got %d", (*res)[1].DocID)
	}

	if (*res)[2].DocID != 5 {
		t.Errorf("third docid Wrong, expected 5, got %d", (*res)[2].DocID)
	}

	if (*res)[3].DocID != 8 {
		t.Errorf("fourth docid Wrong, expected 8, got %d", (*res)[1].DocID)
	}

	arrayEqual = true
	correct = []int64{3, 8, 25, 39, 89, 90}

	if len((*res)[3].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[3].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in fourth posting, expected [1,2,5,6,7,8],got %v", (*res)[3].Pos)
	}
}

func TestPostitionalIntersect(t *testing.T) {
	var list1 = make(PostingList, 0)
	var list2 = make(PostingList, 0)

	list1 = append(list1, Posting{1, []int64{2, 5, 8}})
	list1 = append(list1, Posting{5, []int64{6, 9, 20}})
	list1 = append(list1, Posting{8, []int64{3, 37, 89}})

	list2 = append(list2, Posting{1, []int64{1, 6, 7}})
	list2 = append(list2, Posting{4, []int64{8, 23, 91}})
	list2 = append(list2, Posting{8, []int64{25, 39, 90}})

	var res = list1.Proximity(&list2, 2)

	if len(*res) != 2 {
		t.Errorf("wrong posting count, expected 2, got %d", len(*res))
	}

	if (*res)[0].DocID != 1 {
		t.Errorf("First DocID Wrong, Expected 1, got %d", (*res)[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{1, 2, 5, 6, 7, 8}

	if len((*res)[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in first posting, expected [1,2,5,6,7,8],got %v", (*res)[0].Pos)
	}

	if (*res)[1].DocID != 8 {
		t.Errorf("second docid Wrong, expected 8, got %d", (*res)[1].DocID)
	}

	arrayEqual = true
	correct = []int64{37, 39, 89, 90}

	if len((*res)[1].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[1].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in second posting, expected %v,got %v", correct, (*res)[1].Pos)
	}
}

func TestPhraseIntersect(t *testing.T) {
	var list1 = make(PostingList, 0)
	var list2 = make(PostingList, 0)
	var list3 = make(PostingList, 0)

	list1 = append(list1, Posting{1, []int64{2, 5, 8}})
	list1 = append(list1, Posting{5, []int64{6, 9, 20}})
	list1 = append(list1, Posting{8, []int64{3, 37, 89}})

	list2 = append(list2, Posting{1, []int64{1, 6, 7}})
	list2 = append(list2, Posting{4, []int64{8, 23, 91}})
	list2 = append(list2, Posting{8, []int64{25, 39, 90}})

	list3 = append(list3, Posting{1, []int64{3, 7, 10}})
	list3 = append(list3, Posting{4, []int64{9, 18, 58}})
	list3 = append(list3, Posting{8, []int64{40, 91}})

	var res = list1.PhraseIntersect([]*PostingList{&list2, &list3})

	if len(*res) != 2 {
		t.Errorf("wrong posting count, expected 2, got %d", len(*res))
	}

	if (*res)[0].DocID != 1 {
		t.Errorf("First DocID Wrong, Expected 1, got %d", (*res)[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{5, 6, 7}

	if len((*res)[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in first posting, expected %v,got %v", correct, (*res)[0].Pos)
	}

	if (*res)[1].DocID != 8 {
		t.Errorf("second docid Wrong, expected 8, got %d", (*res)[1].DocID)
	}

	arrayEqual = true
	correct = []int64{89, 90, 91}

	if len((*res)[1].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[1].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in second posting, expected %v,got %v", correct, (*res)[1].Pos)
	}
}

func TestDifference(t *testing.T) {
	var list1 = make(PostingList, 0)
	var list2 = make(PostingList, 0)

	list1 = append(list1, Posting{1, []int64{2, 5, 8}})
	list1 = append(list1, Posting{5, []int64{6, 9, 20}})
	list1 = append(list1, Posting{8, []int64{3, 8, 89}})

	list2 = append(list2, Posting{1, []int64{1, 6, 7}})
	list2 = append(list2, Posting{4, []int64{8, 23, 91}})
	list2 = append(list2, Posting{8, []int64{25, 39, 90}})

	var res = list1.Difference(&list2)

	if len(*res) != 1 {
		t.Errorf("wrong posting count, expected 1, got %d", len(*res))
	}

	if (*res)[0].DocID != 5 {
		t.Errorf("First DocID Wrong, Expected 5, got %d", (*res)[0].DocID)
	}

	var arrayEqual = true
	var correct = []int64{6, 9, 20}

	if len((*res)[0].Pos) != len(correct) {
		arrayEqual = false
	} else {
		for i, v := range (*res)[0].Pos {
			if v != correct[i] {
				arrayEqual = false
			}
		}
	}

	if !arrayEqual {
		t.Errorf("wrong positions in first posting, expected [6 9 20],got %v", (*res)[0].Pos)
	}
}

func TestW2VConstruction(t *testing.T) {
	index, err := BuildIndex("/home/colin/Documents/WPP2/WPP2-INR/docs.txt", "/home/colin/Documents/WPP2/WPP2-INR/binGlove.bin")
	if err != nil {
		t.Errorf("error building index: %v", err)
	}
	t.Logf("%v", index)

	_ = index.EvaluateQuery("diabetes")
}

package spell

import "testing"

func TestJaccard(t *testing.T) {
	a := "bord"
	b := "boardroom"
	var j float32 = 3.0 / 12.0

	if r := Jaccard(a, b, 3, 2); r != j {
		t.Errorf("Expected Jaccard-Index to be %f. Got %f instead", j, r)
	}
}

package spell

import "testing"

func TestBiGram(t *testing.T) {
	const a string = "hallo"
	const k int = 2
	const exLen = 6

	var exGrams = []string{"$h", "ha", "al", "ll", "lo", "o$"}

	if d := ExtractKGrams(a, k); d != nil {
		if len(d) != exLen {
			t.Errorf("Expected K-Gram count of %d. Got %d instead", exLen, len(d))
			return
		}
		for i := range d {
			if d[i] != exGrams[i] {
				t.Errorf("Expected K-Gram %s at Position %d. Got %s instead", exGrams[i], i, d[i])
				return
			}
		}
		return
	}
	t.Errorf("Got nil instead of array of strings")
}

func TestTriGram(t *testing.T) {
	const a string = "eierschalen"
	const k int = 3
	const exLen = 11

	var exGrams = []string{
		"$ei",
		"eie",
		"ier",
		"ers",
		"rsc",
		"sch",
		"cha",
		"hal",
		"ale",
		"len",
		"en$",
	}

	if d := ExtractKGrams(a, k); d != nil {
		if len(d) != exLen {
			t.Errorf("Expected K-Gram count of %d. Got %d instead", exLen, len(d))
			return
		}
		for i := range d {
			if d[i] != exGrams[i] {
				t.Errorf("Expected K-Gram %s at Position %d. Got %s instead", exGrams[i], i, d[i])
				return
			}
		}
		return
	}
	t.Errorf("Got nil instead of array of strings")
}

func TestCountKGrams(t *testing.T) {
	const a string = "eierschalen"
	const k int = 3
	const exLen = 11

	if count := CountKGrams(a, k); count != exLen {
		t.Errorf("Wrong count. Expected %d, got %d", exLen, count)
	}
}

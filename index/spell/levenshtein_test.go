package spell

import "testing"

func TestDistance(t *testing.T) {
	//?????
	const stringA string = "fooo"
	const stringB string = "baro"
	const dist int = 3

	if d := LevenshteinDistance(stringA, stringB); d != dist {
		t.Errorf("Wrong levenshtein-distance expected %d got %d", dist, d)
	}

}

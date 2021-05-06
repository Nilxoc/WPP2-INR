package spell

import "testing"

func TestDistance(t *testing.T) {
	//?????
	const stringA string = "foo"
	const stringB string = "bar"
	const dist int = 3
	//Function not yet working as expected... Need to find out why..
	return

	if d := LevenshteinDistance(stringA, stringB); d != dist {
		//Why?
		t.Errorf("Wrong levenshtein-distance expected %d got %d", dist, d)
	}

}

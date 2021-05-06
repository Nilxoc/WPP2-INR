package query

import (
	"testing"
)

func TestParseAnd(t *testing.T) {
	parseNoError(t, "test1 AND test2")
}

func TestParseProximity(t *testing.T) {
	parseNoError(t, "term4 /3 term5")
}

func TestParsePhrase(t *testing.T) {
	parseNoError(t, "[term1 term2 term3]")
}

func TestParseBrackets(t *testing.T) {
	parseNoError(t, "( term1 AND term2 )")
}

func ParseFull(t *testing.T) {
	parseNoError(t, "(\"term1 term2\" OR term3) AND NOT term4 /3 term5")
}

func parseNoError(t *testing.T, inSample string) {
	q, err := Parse(inSample)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		Print(q)
	}
}

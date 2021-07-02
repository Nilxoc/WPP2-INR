package eval

import (
	"testing"
)

func TestCalcConf(t *testing.T) {
	relList := []int64{1, 3, 5, 6, 7, 9, 10}
	foundList := []int64{1, 2, 5, 9, 4}

	confMat := CalculateConfusion(foundList, relList)

	if !(confMat.tp == 3) {
		t.Errorf("Wrong tp, expected 3, got %d", confMat.tp)
	}

	if !(confMat.fp == 2) {
		t.Errorf("Wrong fp, expected 2, got %d", confMat.fp)
	}

	if !(confMat.fn == 4) {
		t.Errorf("Wrong fn, expected 4, got %d", confMat.fn)
	}
}

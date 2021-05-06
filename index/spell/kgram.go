package spell

//FROM: https://www.forcepoint.com/de/blog/x-labs/simple-n-gram-calculator-pyngram
func ExtractKGrams(token string, k int) []string {
	t := len(token)

	if t == 0 {
		return nil
	}

	arr := make([]string, 0)

	for i := 0; i < (t-k)+1; i++ {
		arr = append(arr, token[i:i+k])
	}

	return arr
}

func CountKGrams(token string, k int) int {
	t := len(token)
	ret := 0

	if t == 0 {
		return 0
	}

	for i := 0; i < (t-k)+1; i++ {
		t += 1
	}

	return ret
}

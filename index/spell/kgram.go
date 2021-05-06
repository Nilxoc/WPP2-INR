package spell

//FROM: https://www.forcepoint.com/de/blog/x-labs/simple-n-gram-calculator-pyngram
func ExtractKGrams(token string, k int) []string {
	s := []rune(token)
	t := len(s)

	if t == 0 {
		return nil
	}

	arr := make([]string, 0)

	for i := 0; i < (t-k)+1; i++ {
		arr = append(arr, token[i:i+k])
	}

	return arr
}

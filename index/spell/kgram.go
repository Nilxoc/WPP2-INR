package spell

//FROM: https://www.forcepoint.com/de/blog/x-labs/simple-n-gram-calculator-pyngram
func ExtractKGrams(token string, k int) []string {
	t := len(token)

	if t == 0 {
		return nil
	}

	arr := make([]string, 0)

	if t < (k - 2) {
		return arr
	}

	arr = append(arr, "$"+token[0:k-1])
	for i := 0; i < (t-k)+1; i++ {
		arr = append(arr, token[i:i+k])
	}
	arr = append(arr, token[len(token)-(k-1):]+"$")

	return arr
}

func CountKGrams(token string, k int) int {
	t_len := len(token) + 2

	if t_len < k {
		return 0
	} else {
		return (t_len - k) + 1
	}
}

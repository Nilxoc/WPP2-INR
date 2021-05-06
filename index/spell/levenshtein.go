package spell

import "fmt"

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func LevenshteinDistance(s string, t string) int {

	//Creating Array's
	m := len(s)
	n := len(t)
	D := make([][]int, m+1)
	for i := 0; i < m+1; i++ {
		D[i] = make([]int, n+1)
	}

	fmt.Println(D)

	//Prepare Matrix
	for i := 0; i <= m; i++ {
		D[i][0] = i
	}
	for j := 0; j <= n; j++ {
		D[0][j] = j
	}

	fmt.Println(D)

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			cost := 0
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			D[i][j] = min(
				D[i-1][j]+1,
				D[i][j-1]+1,
				D[i-1][j-1]+cost,
			)
		}
	}
	return D[m-1][n-1]
}

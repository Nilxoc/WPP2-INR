package spell

import "sort"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func contains(arr []string, v string) bool {
	for _, s := range arr {
		if s == v {
			return true
		}
	}
	return false
}

/* UNUSED
func union(a []string, b []string) []string {
	la := len(a)
	lb := len(b)
	res := make([]string, 0, la+lb)

	for _, e := range a {
		res = append(res, e)
	}

	for _, e := range b {
		if !contains(res, e) {
			res = append(res, e)
		}
	}
	return res
}

func intersect(a []string, b []string) []string {
	la := len(a)
	lb := len(b)
	res := make([]string, 0, min(la, lb))
	for _, e := range a {
		if contains(b, e) {
			res = append(res, e)
		}
	}
	return res
}
*/

func intersectCount(a []string, b []string) float32 {
	res := 0
	for _, e := range a {
		if contains(b, e) {
			res += 1
		}
	}
	return float32(res)
}

func Jaccard(a string, b string, k int) float32 {

	ka := ExtractKGrams(a, k)
	kb := ExtractKGrams(b, k)

	la := float32(len(ka))
	lb := float32(len(kb))

	sort.Strings(ka)
	sort.Strings(kb)

	return intersectCount(ka, kb) / (la + lb - float32(intersectCount(ka, kb)))

}

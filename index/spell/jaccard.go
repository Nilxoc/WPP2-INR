package spell

func Jaccard(a string, b string, matchCount int, k int) float32 {
	la := CountKGrams(a, k)
	lb := CountKGrams(b, k) //Schneller mit dazuspeichern?

	return (float32(matchCount) / float32(la+lb-matchCount))
}

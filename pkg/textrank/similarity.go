package textrank

import "math"

func cosineSimilarity(s1, s2 []string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	counter1 := make(map[string]int)
	counter2 := make(map[string]int)

	for _, word := range s1 {
		counter1[word]++
	}
	for _, word := range s2 {
		counter2[word]++
	}

	dotProduct := 0.0
	for word, count1 := range counter1 {
		dotProduct += float64(count1 * counter2[word])
	}

	magnitudeVec1 := 0.0
	magnitudeVec2 := 0.0
	for _, count := range counter1 {
		magnitudeVec1 += float64(count * count)
	}
	for _, count := range counter2 {
		magnitudeVec2 += float64(count * count)
	}
	magnitudeVec1 = math.Sqrt(magnitudeVec1)
	magnitudeVec2 = math.Sqrt(magnitudeVec2)

	if magnitudeVec1 == 0 || magnitudeVec2 == 0 {
		return 0.0
	}

	return dotProduct / (magnitudeVec1 * magnitudeVec2)
}

func textrankSimilarity(s1, s2 []string) float64 {
	n1 := len(s1)
	n2 := len(s2)

	if n1 <= 1 || n2 <= 1 {
		return 0
	}

	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, word := range s1 {
		set1[word] = true
	}

	for _, word := range s2 {
		set2[word] = true
	}

	common := 0

	for word := range set1 {
		if set2[word] {
			common++
		}
	}
	base := math.Log(float64(n1)) + math.Log(float64(n2))
	return float64(common) / base
}
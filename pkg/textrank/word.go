package textrank

import "sort"

func wordGraph(sents []string, minCount, window, minCooccurrence int) (*SparseMatrix, []string) {
	idxToVocab, vocabToIdx := scanVocabulary(sents, minCount)

	tokens := make([][]string, len(sents))
	for i, sent := range sents {
		tokens[i] = Tokenize(sent)
	}

	matrix := cooccurrence(tokens, vocabToIdx, window, minCooccurrence)

	return matrix, idxToVocab
}

func scanVocabulary(sents []string, minCount int) ([]string, map[string]int) {
	counter := make(map[string]int)
	for _, sent := range sents {
		tokens := Tokenize(sent)
		for _, word := range tokens {
			counter[word]++
		}
	}

	var sortedWords []kv_count
	for word, count := range counter {
		if count >= minCount {
			sortedWords = append(sortedWords, kv_count{word, count})
		}
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].count > sortedWords[j].count
	})

	idxToVocab := make([]string, len(sortedWords))
	vocabToIdx := make(map[string]int)
	for i, kv := range sortedWords {
		idxToVocab[i] = kv.word
		vocabToIdx[kv.word] = i
	}

	return idxToVocab, vocabToIdx
}

func cooccurrence(tokens [][]string, vocabToIdx map[string]int, window, minCooccurrence int) (*SparseMatrix) {
	counter := make(map[int]map[int]int)

	for _, sentTokens := range tokens {
		vocabs := make([]int, 0) 
		for _, token := range sentTokens {
			if idx, ok := vocabToIdx[token]; ok {
				vocabs = append(vocabs, idx)
			}
		}

		var b,e int
		for i, v := range vocabs {
			if window <= 0 {
				b, e = 0, len(vocabs)
			} else {
				b = max(0, i-window)
				e = min(i+window, len(vocabs))
			}

			for j := b; j<e; j++ {
				if i == j {
					continue
				}
				neighbor := vocabs[j]

				if _, ok := counter[v]; !ok {
					counter[v] = make(map[int]int)
				}

				counter[v][neighbor]++
			}
		}
	}

	filteredCounter := make(map[int]map[int]int)
	for token, neighbors := range counter {
		for neighbor, count := range neighbors {
			if count >= minCooccurrence {
				if _, ok := filteredCounter[token]; !ok {
					filteredCounter[token] = make(map[int]int)
				}
				filteredCounter[token][neighbor] = count
			}
		}
		if len(neighbors) == 0 {
			delete(filteredCounter, token)
		}
	}
	row, col, val := CounterToSparseMatrix(filteredCounter)

	matrix := &SparseMatrix{
		row: row,
		col : col,
		val : val,
	}

	return matrix
}


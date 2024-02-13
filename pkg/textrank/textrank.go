package textrank

import (
	"sort"

	"github.com/aglide100/news-scrap/pkg/logger"
)

// reference https://lovit.github.io/nlp/2019/04/30/textrank/

func pagerank(graph *SparseMatrix, df float64, maxIter int) map[int]float64 {
	nor_graph := normalize(graph)

	A := make([]map[int]float64, len(nor_graph.row))
	for i := 0; i < len(nor_graph.row); i++ {
		A[i] = make(map[int]float64)
		for j, col := range nor_graph.col {
			if j < len(nor_graph.row) && nor_graph.row[j] == i {
				A[i][col] = float64(nor_graph.val[j])
			}
		}
	}

	R := make(map[int]float64)
	for i := range nor_graph.row {
		R[i] = 1.0
	}

	bias := (1 - df) * float64(len(nor_graph.row))

	for iter := 0; iter < maxIter; iter++ {
		newR := make(map[int]float64)
		for token, neighbors := range A {
			for neighbor, weight := range neighbors {
				newR[neighbor] += df * weight * R[token]
			}
		}

		for token := range newR {
			newR[token] += bias
		}

		R = newR
	}

	return R
}

func sentGraph(sents []string, minCount int, minSim float64) *SparseMatrix {
	tokens := make([][]string, len(sents))
	for i, sent := range sents {
		tokens[i] = Tokenize(sent)
	}

	rows := []int{}
	cols := []int{}
	data := []float64{}

	for i, tokensI := range tokens {
		for j, tokensJ := range tokens {
			if i >= j {
				continue
			}
			sim := textrankSimilarity(tokensI, tokensJ)
			// sim := cosineSimilarity(tokensI, tokensJ)
			
			if sim < minSim {
				continue
			}
			rows = append(rows, i)
			cols = append(cols, j)
			data = append(data, sim)
		}
	}
	return &SparseMatrix{
		row:   rows,
		col:   cols,
		val: data,
	}
}

func TextRankSentences(sents []string, minCount int, df float64, maxIter, topK int) []*KeySentence {
	g := sentGraph(sents, minCount, 0.3)
	R := pagerank(g, df, maxIter)

	if topK >= len(sents) {
		logger.Info("topK is too big...")
		topK = 3
	}

	scores := make([]kv_score, topK)
	for idx, score := range R {
		scores = append(scores, kv_score{idx, score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	topSent := []*KeySentence{}
	for i := 0; i < topK ; i++ {
		idx := scores[i].idx
		if scores[i].idx != 0 {
			topSent = append(topSent, &KeySentence{idx, R[idx], sents[idx]})
		}
	}

	// sort.Slice(topSent, func(i, j int) bool {
	// 	return topSent[i].Index < topSent[j].Index
	// })

	return topSent
}

func TextRankWordGraph(sents []string, df float64, maxIter, minCount, window, minCooccurrence int) map[string]float64 {
	matrix, idxToVocab := wordGraph(sents, minCount, window, minCooccurrence)

	pageRankScores := pagerank(matrix, df, maxIter)

	textRankScores := make(map[string]float64)
	for idx, score := range pageRankScores {
		word := idxToVocab[idx]
		textRankScores[word] = score
	}

	return textRankScores
}


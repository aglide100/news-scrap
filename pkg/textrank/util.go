package textrank

import (
	"strings"
	"unicode"
)

func normalize(matrix *SparseMatrix) *SparseMatrix {
	rowSums := make([]float64, len(matrix.row))
	val := make([]float64, len(matrix.val))
	copy(val, matrix.val)

	for i := 0; i < len(matrix.row); i++ {
		sum := 0.0
		for j := 0; j < len(matrix.col); j++ {
			idx := i*len(matrix.col) + j
			if idx < len(val) {
				sum += val[idx]
			}
		}
		rowSums[i] = sum
	}

	after := &SparseMatrix{
		row:         make([]int, len(matrix.row)),
		col:         make([]int, len(matrix.col)),
		val: make([]float64, len(matrix.val)),
	}

	copy(after.row, matrix.row)
	copy(after.col, matrix.col)

	for i := 0; i < len(matrix.row); i++ {
		for j := 0; j < len(matrix.col); j++ {
			idx := i*len(matrix.col) + j
			if idx < len(val) {
				after.val[idx] = val[idx] / rowSums[i]
			}
		}
	}
	
	return after
}

func contains(s []string, substr string) bool {
    for _, v := range s {
        if v == substr {
            return true
        }
    }

    return false
}

func Preprocessing(input string) string {
    var result string

    for _, char := range input {
        if unicode.IsGraphic(char) && !unicode.IsSpace(char) && !unicode.IsPunct(char) {
            result += string(char)
        }
    }

    return result
}

func Tokenize(text string) []string {
	tokens := strings.Fields(text)

	var result []string
	for _, token := range tokens {
		current := Preprocessing(token) 

		if len(current) <= 1 || contains(result, current) {
			continue
		}
		
		result = append(result, current)
	}
	return result
}
// func Tokenize(text string) []string {
// 	tokens := strings.Split(text, " ")
	
// 	var result []string
// 	for _, token := range tokens {
// 		if token != "" {
// 			result = append(result, token)
// 		}
// 	}
// 	return result
// }

// func Tokenize(sent string) []string {
// 	n := 3
//     var result []string

//     subword := func(token string, n int) []string {
//         if utf8.RuneCountInString(token) <= n {
//             return []string{token}
//         }
//         var subs []string
//         for i := 0; i<utf8.RuneCountInString(token)-n+1; i++ {
//             subs = append(subs, string([]rune(token)[i:i+n]))
//         }
//         return subs
//     }

//     words := strings.Fields(sent)
//     for _, word := range words {
//         result = append(result, subword(word, n)...)
//     }

//     return result
// }

// func subwordTokenizer(sent string, n int) []string {
//     var result []string

//     subword := func(token string, n int) []string {
//         if len(token) <= n {
//             return []string{token}
//         }
//         var subs []string
//         for i := 0; i < len(token)-n+1; i++ {
//             subs = append(subs, token[i:i+n])
//         }
//         return subs
//     }

//     words := strings.Fields(sent)
//     for _, word := range words {
//         result = append(result, subword(word, n)...)
//     }

//     return result
// }

func CounterToSparseMatrix(counter map[int]map[int]int) ([]int, []int, []float64) {
	row := make([]int, 0)
	col := make([]int, 0)
	val := make([]float64, 0)

	for i, arr := range counter {
		for j, count := range arr {
			row = append(row, i)
			col = append(col, j)
			val = append(val, float64(count))
		}
	}

	return row, col, val
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func SplitText(text string) []string {
	var sentences []string
	var sentence []rune

	for _, char := range text {
		if unicode.IsPunct(char) {
			if char == '.' || char == '?' || char == '!' {
				if len(sentence) > 0 {
					sentences = append(sentences, string(sentence))
					sentence = []rune{}
				}
			}
		} else {
			sentence = append(sentence, char)
		}
	}

	if len(sentence) > 0 {
		sentences = append(sentences, string(sentence))
	}

	return sentences
}
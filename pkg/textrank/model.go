package textrank

type SparseMatrix struct {
	row []int
	col []int
	val []float64
}

type KeySentence struct {
	Index  int
	Score  float64 
	Sentence string 
}

type kv_count struct {
	word  string
	count int
}

type kv_score struct {
	idx   int
	score float64
}
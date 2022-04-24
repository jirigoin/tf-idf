package tfidf

type TF struct {
	TotalWords   int
	WordQuantity int
	Score        float64
}

type IDF struct {
	TotalFiles int
	Df         int
	Score      float64
}

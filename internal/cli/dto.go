package cli

type Flags struct {
	Word string
	File string
}

type Info struct {
	Root        string
	File        string
	Word        string
	Documents   []Document
	FilesNumber int
}

type Result struct {
	Info  Info
	TF    TF
	IDF   IDF
	TFIDF float64
}

type Weigher interface {
	Weigh(info Info) Weigher
}

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

type Document struct {
	Name string
}

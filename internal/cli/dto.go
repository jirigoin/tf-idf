package cli

import (
	"github.com/jirigoin/tf-idf/internal/tfidf"
)

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
	TF    tfidf.TF
	IDF   tfidf.IDF
	TFIDF float64
}

type Weigher interface {
	Weigh(info Info) Weigher
}

type Document struct {
	Name string
}

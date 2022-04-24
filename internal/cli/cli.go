package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/jirigoin/tf-idf/internal/tfidf"
)

const (
	appName     = "Wonderful-TF-IDF-weigher"
	defaultFile = "document_1.txt"
)

func InitFlags(f *Flags) {
	if len(os.Args) > 1 {
		f.Word = os.Args[1]
		flag.CommandLine.StringVar(&f.File, "file", defaultFile, "text file where the word gonna be search.")
		flag.CommandLine.Usage = help
		flag.CommandLine.Parse(os.Args[2:])
	}
	validateFlags(f)
}

func Run(ctx context.Context, flags *Flags) (*Result, error) {
	doneCh := make(chan bool)

	var info Info
	info = info.New(flags)

	var result Result

	go func() {
		result = compute(info)
		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-doneCh:
		return &result, nil
	}
}

func compute(info Info) Result {
	tf := tfidf.TF{}
	idf := tfidf.IDF{}
	var documents []string
	for _, v := range info.Documents {
		documents = append(documents, v.Name)
	}

	r := Result{
		Info: info,
		TF:   tf.Weigh(info.Root, info.File, info.Word),
		IDF:  idf.Weigh(info.Root, info.Word, documents),
	}
	r.TFIDF = r.TF.Score * r.IDF.Score
	return r
}

func validateFlags(f *Flags) {
	if f.Word == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func help() {
	msg := fmt.Sprintf("usage: %s is a simple CLI tool to get the TF-IDF weight of a word in a document\n"+
		"You have to pass a word to be searched\n[OPTIONS]:", appName)
	log.Println(msg)
	flag.PrintDefaults()
}

func (i *Info) New(flag *Flags) Info {
	_, here, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(here), "../../data")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	var documents []Document
	for _, f := range files {
		reg := regexp.MustCompile(`\.txt$`)
		if !f.IsDir() && reg.MatchString(f.Name()) {
			d := Document{Name: f.Name()}
			documents = append(documents, d)
		}
	}
	return Info{
		Root:        dir,
		File:        flag.File,
		Word:        flag.Word,
		Documents:   documents,
		FilesNumber: len(files),
	}
}

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
	"strings"
)

const (
	appName     = "Wonderful-TF-IDF-weigher"
	defaultFile = "document_1.txt"
)

func InitFlags(f *Flags) {
	flag.StringVar(&f.File, "file", defaultFile, "text file where the word gonna be search.")
	flag.Usage = help
	flag.Parse()

	if len(os.Args) > 1 {
		f.Word = os.Args[1]
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
	tf := TF{}
	idf := IDF{}
	r := Result{
		Info: info,
		TF:   tf.Weigh(info),
		IDF:  idf.Weigh(info),
	}
	r.TFIDF = r.TF.Score * r.IDF.Score
	return r
}

func sanitize(text []byte) string {
	ts := strings.ReplaceAll(string(text), "\n", " ")
	ts = strings.ToLower(ts)
	return strings.TrimRight(ts, " ")
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

func stats(searchedWord string, text []byte) (exists bool, times int, totalWords int) {
	tx := sanitize(text)
	w := strings.ToLower(searchedWord)

	times = strings.Count(tx, w)
	if times > 0 {
		exists = true
	}
	totalWords = len(strings.Split(tx, " "))

	return
}

func (i *Info) New(flag *Flags) Info {
	_, here, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(here), "../..")
	if flag.File == defaultFile {
		dir = filepath.Join(dir, "/data")
	}
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

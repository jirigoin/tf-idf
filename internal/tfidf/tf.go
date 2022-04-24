package tfidf

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (t TF) Weigh(directory, file, word string) TF {
	dir := filepath.Join(directory, file)
	fi, err := os.OpenFile(dir, os.O_RDONLY, 0o111)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(fi)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		t.getStats(word, buf[:n])
	}
	t.tf()
	return t
}

func (t *TF) getStats(searchedWord string, text []byte) {
	_, times, words := stats(searchedWord, text)
	t.WordQuantity += times
	t.TotalWords += words
}

func (t *TF) tf() {
	t.Score = float64(t.WordQuantity) / float64(t.TotalWords)
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

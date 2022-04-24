package tfidf

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (idf IDF) Weigh(directory, word string, documents []string) IDF {
	idf.TotalFiles = len(documents)

	for _, d := range documents {
		if idf.existsInDocument(directory, word, d) {
			idf.Df++
		}
	}
	var x float64
	if idf.Df != 0 {
		x = math.Log10(float64(idf.TotalFiles) / float64(idf.Df))
	}
	result, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", x), 64)
	idf.Score = result

	return idf
}

func (idf *IDF) existsInDocument(root string, word string, document string) bool {
	dir := filepath.Join(root, document)
	fi, err := os.OpenFile(dir, os.O_RDONLY, 0o111)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var exists bool

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
		exists = exists || idf.search(word, buf[:n])
	}

	return exists
}

func (idf *IDF) search(searchedWord string, text []byte) bool {
	exists, _, _ := stats(searchedWord, text)
	return exists
}

func sanitize(text []byte) string {
	ts := strings.ReplaceAll(string(text), "\n", " ")
	ts = strings.ToLower(ts)
	return strings.TrimRight(ts, " ")
}

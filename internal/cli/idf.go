package cli

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func (idf IDF) Weigh(info Info) IDF {
	idf.TotalFiles = len(info.Documents)

	for _, d := range info.Documents {
		if idf.existsInDocument(info, d) {
			idf.Df++
		}
	}
	result, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", math.Log10(float64(idf.TotalFiles)/float64(idf.Df))), 64)
	idf.Score = result

	return idf
}

func (idf *IDF) existsInDocument(info Info, document Document) bool {
	dir := filepath.Join(info.Root, document.Name)
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
		exists = exists || idf.search(info.Word, buf[:n])
	}

	return exists
}

func (idf *IDF) search(searchedWord string, text []byte) bool {
	exists, _, _ := stats(searchedWord, text)
	return exists
}

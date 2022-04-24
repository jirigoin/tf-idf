package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jirigoin/tf-idf/internal/cli"
)

func main() {
	n := make(chan os.Signal)

	signal.Notify(n, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-n
		cancel()
	}()

	flags := &cli.Flags{}

	cli.InitFlags(flags)

	r, err := cli.Run(ctx, flags)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Search: \"%s\" in the file %s. Appears: %d times\n", flags.Word, flags.File, r.TF.WordQuantity)

	fmt.Printf("tf (\"%s\", %s) = %.3f\n", r.Info.Word, r.Info.File, r.TF.Score)

	fmt.Printf("idf (\"%s\", appears in= %d of %d files) = %.3f\n", r.Info.Word, r.IDF.Df, r.IDF.TotalFiles, r.IDF.Score)

	fmt.Println("--------------------------")
	fmt.Printf("tfidf = %.3f\n", r.TFIDF)
}

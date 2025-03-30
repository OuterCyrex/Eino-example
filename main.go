package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	uuid2 "github.com/google/uuid"
	"io"
)

const (
	prefix = "OuterCyrex:"
	index  = "OuterIndex"
)

func main() {
	ctx := context.Background()

	r, err := InitRAGEngine(ctx, index, prefix)
	if err != nil {
		panic(err)
	}

	doc, err := r.Loader.Load(ctx, document.Source{
		URI: "./test_txt/mysql-1.md",
	})
	if err != nil {
		panic(err)
	}

	docs, err := r.Splitter.Transform(ctx, doc)
	if err != nil {
		panic(err)
	}

	for _, d := range docs {
		uuid, _ := uuid2.NewUUID()
		d.ID = uuid.String()
	}

	err = r.InitVectorIndex(ctx)
	if err != nil {
		panic(err)
	}

	_, err = r.Indexer.Store(ctx, docs)
	if err != nil {
		panic(err)
	}

	var query string

	for {
		_, _ = fmt.Scan(&query)
		output, err := r.Generate(ctx, query)
		if err != nil {
			panic(err)
		}
		for {
			o, err := output.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(o.Content)
		}
	}
}

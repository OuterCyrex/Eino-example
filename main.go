package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	uuid2 "github.com/google/uuid"
	"io"
)

var (
	indexName = "index_test"
	prefix    = "2025:"
)

func main() {
	ctx := context.Background()

	rag, err := InitRAGClient(ctx, prefix, indexName)
	if err != nil {
		panic(err)
	}

	docs, err := rag.Loader.Load(ctx, document.Source{
		URI: "./test_txt/mysql-1.md",
	})
	if err != nil {
		panic(err)
	}

	ted, err := rag.Transformer.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}
	for _, doc := range ted {
		uuid, _ := uuid2.NewUUID()
		doc.ID = uuid.String()
	}

	err = rag.InitVectorIndex(ctx)
	if err != nil {
		panic(err)
	}

	_, err = rag.Indexer.Store(ctx, ted)
	if err != nil {
		panic(err)
	}

	output, err := rag.Generate(ctx, "介绍一下什么是存储引擎呢")
	if err != nil {
		panic(err)
	}

	for {
		o, e := output.Recv()
		if e == io.EOF {
			break
		}
		if e != nil {
			panic(e)
		}
		fmt.Println(o.Content)
	}
}

package main

import (
	"context"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
)

func (r *RAGEngine) newRetriever(ctx context.Context) {
	re, err := redisRet.NewRetriever(ctx, &redisRet.RetrieverConfig{
		Client:       r.redis,
		Index:        r.indexName,
		VectorField:  "vector_content",
		Dialect:      2,
		ReturnFields: []string{"vector_content", "content"},
		TopK:         1,
		Embedding:    r.embedder,
	})
	if err != nil {
		r.Err = err
	}
	r.Retriever = re
}

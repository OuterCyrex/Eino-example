package main

import (
	"context"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
)

func (r *RAGEngine) newIndexer(ctx context.Context) {
	i, err := redisInd.NewIndexer(ctx, &redisInd.IndexerConfig{
		Client:           r.redis,
		KeyPrefix:        r.prefix,
		DocumentToHashes: nil,
		BatchSize:        10,
		Embedding:        r.embedder,
	})
	if err != nil {
		r.Err = err
	}
	r.Indexer = i
}

func (r *RAGEngine) InitVectorIndex(ctx context.Context) error {
	if _, err := r.redis.Do(ctx, "FT.INFO", r.indexName).Result(); err == nil {
		return nil
	}

	createIndexArgs := []interface{}{
		"FT.CREATE", r.indexName,
		"ON", "HASH",
		"PREFIX", "1", r.prefix,
		"SCHEMA",
		"content", "TEXT",
		"vector_content", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", r.dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	if err := r.redis.Do(ctx, createIndexArgs...).Err(); err != nil {
		return err
	}

	if _, err := r.redis.Do(ctx, "FT.INFO", r.indexName).Result(); err != nil {
		return err
	}
	return nil
}

package main

import (
	"context"
	"fmt"
	redisIndexer "github.com/cloudwego/eino-ext/components/indexer/redis"
)

func (r *RAGClient) newIndexer(ctx context.Context) {
	i, err := redisIndexer.NewIndexer(ctx, &redisIndexer.IndexerConfig{
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

func (r *RAGClient) InitVectorIndex(ctx context.Context) error {
	_, err := r.redis.Do(ctx, "FT.INFO", r.indexName).Result()
	if err == nil {
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
		return fmt.Errorf("failed to create index: %w", err)
	}

	_, err = r.redis.Do(ctx, "FT.INFO", r.indexName).Result()
	if err != nil {
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	return nil
}

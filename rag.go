package main

import (
	"context"
	"eino/config"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	embedding "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino-ext/components/model/ark"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

type RAGEngine struct {
	indexName string
	prefix    string
	config    *config.ParamsConfig
	dimension int

	redis    *redis.Client
	embedder *embedding.Embedder

	Err error

	Loader    *file.FileLoader
	Splitter  document.Transformer
	Retriever *redisRet.Retriever
	Indexer   *redisInd.Indexer
	ChatModel *ark.ChatModel
}

func InitRAGEngine(ctx context.Context, index string, prefix string) (*RAGEngine, error) {
	r, err := initRAGEngine(ctx, index, prefix)
	if err != nil {
		return nil, err
	}

	r.newLoader(ctx)
	r.newSplitter(ctx)
	r.newIndexer(ctx)
	r.newRetriever(ctx)
	r.newChatModel(ctx)

	return r, nil
}

func initRAGEngine(ctx context.Context, index string, prefix string) (*RAGEngine, error) {

	c := config.Map()

	embedder, err := embedding.NewEmbedder(ctx, &embedding.EmbeddingConfig{
		APIKey: c.ApiKey,
		Model:  c.Embedding,
	})

	if err != nil {
		return nil, err
	}

	return &RAGEngine{
		indexName: index,
		prefix:    prefix,
		config:    c,
		dimension: 4096,

		redis: redis.NewClient(&redis.Options{
			Addr:          fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
			Protocol:      2,
			UnstableResp3: true,
		}),
		embedder: embedder,

		Loader:    nil,
		Splitter:  nil,
		Retriever: nil,
		Indexer:   nil,
		ChatModel: nil,
	}, nil
}

var systemPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

here's documents searched for you:
==== doc start ====
	  {documents}
==== doc end ====
`

func (r *RAGEngine) Generate(ctx context.Context, query string) (*schema.StreamReader[*schema.Message], error) {
	docs, err := r.Retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, err
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(docs)
	fmt.Println("-------------------------------------------")

	tpl := prompt.FromMessages(schema.FString, []schema.MessagesTemplate{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("question: {content}"),
	}...)

	messages, err := tpl.Format(ctx, map[string]any{
		"documents": docs,
		"content":   query,
	})
	if err != nil {
		return nil, err
	}

	return r.ChatModel.Stream(ctx, messages)
}

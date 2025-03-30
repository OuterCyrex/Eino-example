package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
)

func (r *RAGClient) newChatModel(ctx context.Context) {
	c, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: r.config.ApiKey,
		Model:  r.config.ChatModel,
	})

	if err != nil {
		r.Err = err
	}

	r.ChatModel = c
}

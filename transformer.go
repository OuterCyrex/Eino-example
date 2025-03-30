package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
)

func (r *RAGEngine) newTransformer(ctx context.Context) {
	config := &markdown.HeaderConfig{
		Headers: map[string]string{
			"#": "title",
		},
		TrimHeaders: false}
	t, err := markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		r.Err = err
	}
	r.Transformer = t
}

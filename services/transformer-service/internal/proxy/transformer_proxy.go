// Package proxy
package proxy

import (
	"context"
	"fmt"

	pbTransformer "github.com/neihtq/tap-lingo/gen/go/proto/transformer/v1"
	"github.com/neihtq/tap-lingo/internal/transformer"

	pbTokenizer "github.com/neihtq/tap-lingo/gen/go/tokenizer/v1"
)

type ProxyServer interface {
	Transform(ctx context.Context, req *pbTransformer.TransformRequest) (*pbTransformer.TransformResponse, error)
}

type TransformerProxy struct {
	pbTransformer.UnimplementedTransformerServiceServer

	Transformer            transformer.Transformer
	TokenizerServiceClient pbTokenizer.TokenizerServiceClient
}

func (tp *TransformerProxy) Transform(ctx context.Context, req *pbTransformer.TransformRequest) (*pbTransformer.TransformResponse, error) {
	url := req.GetUrl()
	transformResult := tp.Transformer.TransformArticle(url)
	if transformResult.Result == transformer.Fail {
		fmt.Printf("[WARN] Cannot transform %s \n", url)
	}

	tokenizerRequest := pbTokenizer.TokenizeRequest{Text: transformResult.PlainText}
	res, err := tp.TokenizerServiceClient.Tokenize(ctx, &tokenizerRequest)
	if err != nil {
		fmt.Printf("[ERROR] TokenizerServicer returned error %v \n", err)
		return &pbTransformer.TransformResponse{
			Title:       transformResult.Title,
			Byline:      transformResult.Byline,
			ImageUrl:    transformResult.ImageURL,
			HtmlContent: transformResult.HTML,
			PlaintText:  transformResult.PlainText,
		}, nil

	}

	tokens := make([]*pbTransformer.Token, len(res.Tokens))
	for i, t := range res.Tokens {
		tokens[i] = &pbTransformer.Token{
			Token: t.Token,
			Start: t.Start,
			End:   t.End,
		}
	}

	response := &pbTransformer.TransformResponse{
		Title:       transformResult.Title,
		Byline:      transformResult.Byline,
		ImageUrl:    transformResult.ImageURL,
		HtmlContent: transformResult.HTML,
		PlaintText:  transformResult.PlainText,
		Tokens:      tokens,
	}
	return response, nil
}

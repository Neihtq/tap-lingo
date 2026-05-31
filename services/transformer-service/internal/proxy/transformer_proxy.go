// Package proxy
package proxy

import (
	"context"
	"fmt"

	pb "github.com/neihtq/tap-lingo/gen/go/proto/transformer/v1"
	"github.com/neihtq/tap-lingo/services/transformer-service/internal/transformer"
)

type ProxyServer interface {
	Transform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error)
}

type TransformerProxy struct {
	pb.UnimplementedTransformerServiceServer

	Transformer transformer.Transformer
}

func (s *TransformerProxy) Transform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	url := req.GetUrl()
	transformResult := s.Transformer.TransformArticle(url)
	if transformResult.Result == transformer.Fail {
		fmt.Println("[WARN] Cannot transform %s", url)
	}

	nGrams := []string{"a", "b", "c"}
	response := &pb.TransformResponse{
		Title:       transformResult.Title,
		Byline:      transformResult.Byline,
		ImageUrl:    transformResult.ImageURL,
		HtmlContent: transformResult.HTML,
		PlaintText:  transformResult.PlainText,
		NGrams:      nGrams,
	}

	return response, nil
}

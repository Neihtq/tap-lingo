// Package proxy
package proxy

import (
	"context"
	"fmt"

	pb "github.com/neihtq/tap-lingo/gen/go/proto/transformer/v1"
)

type Server struct {
	pb.UnimplementedTransformerServiceServer
}

func (s *Server) Transform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	url := req.GetUrl()
	fmt.Println(url)

	content := "SomeContent"
	nGrams := []string{"a", "b", "c"}
	response := &pb.TransformResponse{
		Content: content,
		NGrams:  nGrams,
	}

	return response, nil
}

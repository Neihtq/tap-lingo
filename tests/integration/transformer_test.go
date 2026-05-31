package integration

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/neihtq/tap-lingo/gen/go/proto/transformer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "localhost:50051"

func TestTransform(t *testing.T) {
	url := "www.google.com"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("fail to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTransformerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := pb.TransformRequest{Url: url}
	res, err := client.Transform(ctx, &request)
	if err != nil {
		t.Errorf("transform return failure: %v", err)
	}
	log.Printf("Response: %s", res.GetContent())
}

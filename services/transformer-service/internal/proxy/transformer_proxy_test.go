package proxy

import (
	"context"
	"errors"
	"testing"

	"github.com/neihtq/tap-lingo/internal/transformer"
	"google.golang.org/grpc"

	pbTransformer "github.com/neihtq/tap-lingo/gen/go/proto/transformer/v1"
	pbTokenizer "github.com/neihtq/tap-lingo/gen/go/tokenizer/v1"
)

type mockTransformer struct {
	onTransformArticle func(url string) transformer.TransformResult
}

func (m *mockTransformer) TransformArticle(url string) transformer.TransformResult {
	return m.onTransformArticle(url)
}

type mockTokenizerClient struct {
	pbTokenizer.TokenizerServiceClient
	onTokenize func(ctx context.Context, in *pbTokenizer.TokenizeRequest) (*pbTokenizer.TokenizeResponse, error)
}

func (m *mockTokenizerClient) Tokenize(ctx context.Context, in *pbTokenizer.TokenizeRequest, opts ...grpc.CallOption) (*pbTokenizer.TokenizeResponse, error) {
	return m.onTokenize(ctx, in)
}

func TestTransform(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Transformation and Teoknization", func(t *testing.T) {
		mockTransformer := &mockTransformer{
			onTransformArticle: func(url string) transformer.TransformResult {
				return transformer.TransformResult{
					Result:    transformer.Success,
					Title:     "Test Title",
					PlainText: "Test Text",
				}
			},
		}

		mockTokenizer := &mockTokenizerClient{
			onTokenize: func(ctx context.Context, in *pbTokenizer.TokenizeRequest) (*pbTokenizer.TokenizeResponse, error) {
				if in.Text != "Test Text" {
					t.Errorf("Expected tokenization text 'Test Text', but was '%s'", in.Text)
				}
				return &pbTokenizer.TokenizeResponse{
					Tokens: []*pbTokenizer.Token{
						{Token: "Test", Start: 0, End: 4},
						{Token: "Text", Start: 6, End: 10},
					},
				}, nil
			},
		}

		server := &TransformerProxy{
			Transformer:            mockTransformer,
			TokenizerServiceClient: mockTokenizer,
		}

		res, err := server.Transform(ctx, &pbTransformer.TransformRequest{Url: "https://test.com"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got %s", res.Title)
		}
		if len(res.Tokens) != 2 {
			t.Fatalf("Expected 2 tokens, got %d", len(res.Tokens))
		}
		if res.Tokens[0].Token != "Test" {
			t.Errorf("Expected first token to be 'Test', got '%s'", res.Tokens[0].Token)
		}
	})

	t.Run("Tokenizer Fails - Returns Partial Response Without Tokens", func(t *testing.T) {
		mockTransformer := &mockTransformer{
			onTransformArticle: func(url string) transformer.TransformResult {
				return transformer.TransformResult{
					Result:    transformer.Success,
					Title:     "Test Title",
					PlainText: "Test Text",
				}
			},
		}

		mockTokenizer := &mockTokenizerClient{
			onTokenize: func(ctx context.Context, in *pbTokenizer.TokenizeRequest) (*pbTokenizer.TokenizeResponse, error) {
				return nil, errors.New("gRPC service unavailable")
			},
		}

		server := &TransformerProxy{
			Transformer:            mockTransformer,
			TokenizerServiceClient: mockTokenizer,
		}
		res, err := server.Transform(ctx, &pbTransformer.TransformRequest{Url: "https://test.com"})
		if err != nil {
			t.Fatalf("Expected nil error despite tokenixer failure, got %v", err)
		}
		if res == nil {
			t.Fatal("Expected a non-nil response")
		}
		if res.Title != "Test Title" {
			t.Errorf("Expected tite to be 'Test Title', got '%s'", res.Title)
		}
		if len(res.Tokens) != 0 {
			t.Errorf("Expected 0 tokens, got %d", len(res.Tokens))
		}
	})
}

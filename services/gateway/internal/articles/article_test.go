package articles

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pbTransformer "github.com/neihtq/tap-lingo/transformer/gen/go/proto/transformer/v1"
	"google.golang.org/grpc"
)

type mockTransformerClient struct {
	TransformFunc func() (*pbTransformer.TransformResponse, error)
}

func (m *mockTransformerClient) Transform(ctx context.Context, in *pbTransformer.TransformRequest, opts ...grpc.CallOption) (*pbTransformer.TransformResponse, error) {
	if m.TransformFunc == nil {
		panic("TransformFunc was not configured in the test")
	}

	return m.TransformFunc()
}

func transformFuncSucceeds() (*pbTransformer.TransformResponse, error) {
	return &pbTransformer.TransformResponse{
		Title:       "Transform Test",
		Byline:      "John Doe",
		ImageUrl:    "https://example.com/image.png",
		HtmlContent: "<h1>Test</h1>",
		PlaintText:  "Test",
	}, nil
}

func transformFuncFails() (*pbTransformer.TransformResponse, error) {
	return nil, errors.New("Transform failed")
}

func TestPostArticle(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(m *mockTransformerClient)
		expectedStatus int
		validateRes    func(t *testing.T, body string)
	}{
		{
			name:           "Valid Request (201 Created)",
			requestBody:    `{"url": "https://test.article.local"}`,
			mockSetup:      func(m *mockTransformerClient) { m.TransformFunc = transformFuncSucceeds },
			expectedStatus: http.StatusCreated,
			validateRes: func(t *testing.T, body string) {
				var resp PostArticleResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Title != "Transform Test" || resp.Byline != "John Doe" {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name:           "Transformer Service Fails (500 Internal Server error)",
			requestBody:    `{"url": "https://test.article.local"}`,
			mockSetup:      func(m *mockTransformerClient) { m.TransformFunc = transformFuncFails },
			expectedStatus: http.StatusInternalServerError,
			validateRes: func(t *testing.T, body string) {
				if !strings.Contains(body, "Transform failed.") {
					t.Errorf("expected 'Transform failed.' error message, but got: %s", body)
				}
			},
		},
		{
			name:           "Missing URL Field (422 Unprocessable)",
			requestBody:    `{"url": ""}`,
			expectedStatus: http.StatusUnprocessableEntity,
			mockSetup:      func(m *mockTransformerClient) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockTransformerClient{}
			tt.mockSetup(mockClient)
			handler := &ArticlesHandlerImpl{transformerclient: mockClient}

			req := httptest.NewRequest("POST", "/articles", strings.NewReader(tt.requestBody))
			rr := httptest.NewRecorder()

			handler.HandlePostArticle(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.validateRes != nil {
				tt.validateRes(t, rr.Body.String())
			}
		})
	}
}

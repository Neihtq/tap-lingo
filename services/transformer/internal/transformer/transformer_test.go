package transformer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransformArticle(t *testing.T) {
	// Mock HTTP server with sample HTML
	networkCallCount := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<!DOCTYPE html>
			<html>
			<head>
				<title>Test Article Title</title>
				<meta name="author" content="Jane Doe">
			</head>
			<body>
				<article>
					<h1>Test Article Title</title>
					<p>This is the core content of the article</p>
				</article>
			</body>
			</html>`)
		networkCallCount++
	}))
	defer mockServer.Close()

	articleTransformer := NewArticleTransformer()
	t.Run("Cache Miss - successful Fetch and Parse", func(t *testing.T) {
		res := articleTransformer.TransformArticle(mockServer.URL)

		if res.Result != Success {
			t.Fatalf("Expected Success, go %s", res.Result)
		}
		if res.Title != "Test Article Title" {
			t.Errorf("Expected Title 'Test Article Title', got '%s'", res.Title)
		}
		if res.Byline != "Jane Doe" {
			t.Errorf("Expected Byline 'Jane Doe', got '%s'", res.Byline)
		}
		if res.HTML == "" || res.PlainText == "" {
			t.Error("Expected HTML and PlainText to be populated, but they were empty")
		}
		if networkCallCount != 1 {
			t.Errorf("Expected 1 network call, got %d", networkCallCount)
		}
	})

	t.Run("Cache Hit - Does Not Call Network", func(t *testing.T) {
		res := articleTransformer.TransformArticle(mockServer.URL)
		if res.Result != Success {
			t.Fatalf("Expected Success on cache hit, got %s", res.Result)
		}
		if networkCallCount != 1 {
			t.Errorf("Expected network call count to stay 1, but it was %d", networkCallCount)
		}
	})

	t.Run("Fetch Failure", func(t *testing.T) {
		res := articleTransformer.TransformArticle("http://fake-url-tap-lingo.local")

		if res.Result != Fail {
			t.Errorf("Expected Fail result, got %s", res.Result)
		}
	})
}

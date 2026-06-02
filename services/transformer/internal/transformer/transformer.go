// Package transformer
package transformer

import (
	"bytes"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Result string

const (
	Success Result = "Success"
	Fail    Result = "Fail"
)

type TransformResult struct {
	Result    Result
	Title     string
	Byline    string
	ImageURL  string
	HTML      string
	PlainText string
}

type Transformer interface {
	TransformArticle(url string) TransformResult
}

type ArticleTransformer struct {
	Cache *expirable.LRU[string, readability.Article]
}

func NewArticleTransformer() *ArticleTransformer {
	return &ArticleTransformer{
		Cache: expirable.NewLRU[string, readability.Article](5, nil, time.Minute*10),
	}
}

func (a *ArticleTransformer) TransformArticle(url string) TransformResult {
	var article readability.Article
	if val, ok := a.Cache.Get(url); ok {
		article = val
	} else {
		readable, err := readability.FromURL(url, 30*time.Second)
		if err != nil {
			return TransformResult{Result: Fail}
		}
		article = readable
		a.Cache.Add(url, article)
	}

	var htmlBuf, textBuf bytes.Buffer
	if err := article.RenderHTML(&htmlBuf); err != nil {
		return TransformResult{Result: Fail}
	}
	if err := article.RenderText(&textBuf); err != nil {
		return TransformResult{Result: Fail}
	}

	return TransformResult{
		Result:    Success,
		Title:     article.Title(),
		Byline:    article.Byline(),
		ImageURL:  article.ImageURL(),
		HTML:      htmlBuf.String(),
		PlainText: textBuf.String(),
	}
}

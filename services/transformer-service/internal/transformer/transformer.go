// Package transformer
package transformer

import (
	"bytes"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
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
	TransformArticle(url string) (TransformResult, error)
}

type ArticleTransformer struct{}

func (a *ArticleTransformer) TransformArticle(url string) TransformResult {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return TransformResult{Result: Fail}
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

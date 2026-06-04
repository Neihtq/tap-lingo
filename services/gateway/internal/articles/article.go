// Package articles
package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pbTransformer "github.com/neihtq/tap-lingo/clients/transformer/proto/transformer/v1"
)

type PostArticleRequest struct {
	URL string `json:"url"`
}

type PostArticleResponse struct {
	Title     string `json:"title"`
	Byline    string `json:"byline"`
	ImageURL  string `json:"url"`
	HTML      string `json:"html"`
	PlainText string `json:"plainText"`
}

type ArticlesHandler interface {
	HandlePostArticle(w http.ResponseWriter, r *http.Request)
}

type ArticlesHandlerImpl struct {
	transformerclient pbTransformer.TransformerServiceClient
}

func (ah *ArticlesHandlerImpl) HandlePostArticle(w http.ResponseWriter, r *http.Request) {
	var req PostArticleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Malformed JSON payload", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Missing required field 'URL'", http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("Received POST article request with url '%s'\n", req.URL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	transformReq := pbTransformer.TransformRequest{Url: req.URL}
	transformRes, err := ah.transformerclient.Transform(ctx, &transformReq)
	if err != nil {
		fmt.Printf("[ERROR] Transformer failed for URL '%s'.\nError: %v\n", req.URL, err)
		http.Error(w, "Transform failed.", http.StatusInternalServerError)
		return
	}

	resp := PostArticleResponse{
		Title:     transformRes.GetTitle(),
		Byline:    transformRes.GetByline(),
		ImageURL:  transformRes.GetImageUrl(),
		HTML:      transformRes.GetHtmlContent(),
		PlainText: transformRes.GetPlaintText(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

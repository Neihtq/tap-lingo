// Package gateway
package gateway

import (
	"net/http"
	"sync"
)

const PostArticlesEndpoint = "POST /articles"

type GatewayService struct {
	mu   sync.Mutex
	Port string
}

func (gs *GatewayService) StartService() {
	http.NewServeMux()
}

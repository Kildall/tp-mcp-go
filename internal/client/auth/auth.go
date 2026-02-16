package auth

import "net/http"

// Strategy defines how authentication is applied to HTTP requests
type Strategy interface {
	ApplyAuth(req *http.Request)
}

type accessTokenStrategy struct {
	token string
}

// NewAccessTokenStrategy creates a Strategy that adds ?access_token= query parameter
func NewAccessTokenStrategy(token string) Strategy {
	return &accessTokenStrategy{token: token}
}

func (s *accessTokenStrategy) ApplyAuth(req *http.Request) {
	q := req.URL.Query()
	q.Set("access_token", s.token)
	req.URL.RawQuery = q.Encode()
}

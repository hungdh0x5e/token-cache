package tokencache

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewClientCredentialGetter(config clientcredentials.Config) TokenGetter {
	return ClientCredentialGetter{config: config}
}

type ClientCredentialGetter struct {
	config clientcredentials.Config
}

// FetchToken implement TokenGetter interface
func (o ClientCredentialGetter) FetchToken(ctx context.Context) (*oauth2.Token, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	return o.config.Token(ctx)
}

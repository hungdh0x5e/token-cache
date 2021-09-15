package tokencache

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewClientCredentialGetter(config clientcredentials.Config, httpClient *http.Client) TokenGetter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return ClientCredentialGetter{config: config, httpClient: httpClient}
}

type ClientCredentialGetter struct {
	config     clientcredentials.Config
	httpClient *http.Client
}

// FetchToken implement TokenGetter interface
func (o ClientCredentialGetter) FetchToken(ctx context.Context) (*oauth2.Token, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, o.httpClient)

	return o.config.Token(ctx)
}

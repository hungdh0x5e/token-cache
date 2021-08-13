package tokencache_test

import (
	"context"
	"fmt"
	"time"

	tokencache "github.com/hungdh0x5e/token-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Step 1
type CustomTokenGetter struct{}

// FetchToken implement TokenGetter interface
func (o CustomTokenGetter) FetchToken(ctx context.Context) (*oauth2.Token, error) {
	//Request authenticate
	//Parse response into oauth2.Token
	//Return token, err
	return nil, nil
}

// Step 2
func NewClient(tokenCache *tokencache.TokenCache) Client {
	return Client{cacheToken: tokenCache}
}

type Client struct {
	cacheToken *tokencache.TokenCache
}

func (c Client) UserInfo(userID string) {
	token, err := c.cacheToken.GetToken(context.Background())
	if err != nil {
		fmt.Printf("UserInfo: %v, err: %v\n", userID, err)
		return
	}
	fmt.Printf("UserInfo: %v, token: %v\n", userID, token.AccessToken)
}

func Example() {
	tokenGetter := CustomTokenGetter{}
	client := NewClient(tokencache.NewTokenCache(tokenGetter))

	fmt.Println("Must same access token for all go-routine")
	for i := 1; i < 5; i++ {
		go client.UserInfo(fmt.Sprintf("case_%.2d", i))
		//fmt.Println("==================================")
		//time.Sleep(time.Duration(i) * time.Second)
	}

	time.Sleep(5 * time.Second)
}

func ExampleClientCredentialGetter() {
	tokenGetter := tokencache.NewClientCredentialGetter(
		clientcredentials.Config{
			ClientID:     "REPLACE_CLIENT_ID",
			ClientSecret: "REPLACE_CLIENT_SECRET",
			Scopes:       []string{""},
			TokenURL:     "https://domain.com/oauth/token",
		},
	)

	client := NewClient(tokencache.NewTokenCache(tokenGetter))

	fmt.Println("Must same access token for all go-routine")
	for i := 1; i < 5; i++ {
		go client.UserInfo(fmt.Sprintf("case_%.2d", i))
		//fmt.Println("==================================")
		//time.Sleep(time.Duration(i) * time.Second)
	}

	time.Sleep(5 * time.Second)
}

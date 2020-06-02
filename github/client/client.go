// Package client provides helpers to initialize a GitHub client.
package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/die-net/lrucache"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

var ghcInstance *github.Client
var ghcOnce sync.Once

var ghcCachingInstance *github.Client
var ghcCachingOnce sync.Once

// Singleton returns a GitHub client singleton.
//
// A GitHub personal access token is required.
//
// Singleton will try to read the token from the environment variable
// GITHUB_TOKEN or read it from the operating system keychain.
//
// To add the token to the macOS keychain you can use the command line
// utility "security" like this:
//
//   security add-generic-password -a github -s GITHUB_TOKEN -w
//
// To add the token to GNOME keyring use "secret-tool":
//
//   secret-tool store --label="GitHub Token" service GITHUB_TOKEN username github
func Singleton() (*github.Client, error) {
	var err error
	var creds string

	ghcOnce.Do(func() {
		creds, err = getCreds()
		if err == nil {
			ghcInstance, err = newGHClientFromToken(creds)
		}
	})

	return ghcInstance, err
}

// CachingSingleton similar to Singleton but with HTTP caching enabled.
//
// Supported cache URLs:
//   * file:///path/to/cache/dir (disk cache)
//   * mem: (memory cache)
func CachingSingleton(url string) (*github.Client, error) {
	var err error
	var creds string
	var cache httpcache.Cache

	ghcCachingOnce.Do(func() {
		creds, err = getCreds()
		if err == nil {
			cache, err = cacheForURL(url)
			if err == nil {
				ghcCachingInstance, err = newCachingGHClientFromToken(creds, cache)
			}
		}
	})

	return ghcCachingInstance, err
}

func cacheForURL(u string) (httpcache.Cache, error) {
	pURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	switch pURL.Scheme {
	case "mem":
		return httpcache.NewMemoryCache(), nil
	case "file":
		return diskcache.New(pURL.Path), nil
	case "lru":
		return lrucache.New(1048576, 3600), nil
	default:
		return nil, fmt.Errorf("Cache type not supported")
	}
}

func getCreds() (string, error) {
	creds := os.Getenv("GITHUB_TOKEN")
	if creds != "" {
		return creds, nil
	}

	creds, err := keyring.Get("GITHUB_TOKEN", "github")
	if err != nil {
		return "", fmt.Errorf("GitHub token not found in keyring. E: %v", err)
	}

	return creds, nil

}

func newGHClientFromToken(token string) (*github.Client, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token can't be empty")
	}
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func newCachingGHClientFromToken(token string, cache httpcache.Cache) (*github.Client, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token can't be empty")
	}

	oauthTransport := &oauth2.Transport{
		Source: oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		),
	}
	ct := httpcache.NewTransport(cache)
	ct.Transport = oauthTransport

	httpClient := &http.Client{
		Transport: ct,
	}

	return github.NewClient(httpClient), nil
}

func newGHClientFromFile(creds string) *github.Client {
	ctx := context.Background()

	key, err := ioutil.ReadFile(creds)
	if err != nil {
		panic(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: strings.TrimSpace(string(key))},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

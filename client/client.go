package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/patricekahlhorn/xdcc-search/search"
	"github.com/patricekahlhorn/xdcc-search/search/provider"
	"strings"
	"time"
)

type Client struct {
	timeout   time.Duration
	providers []provider.Provider
}

type Option func(*Client)

func NewClient(options ...Option) *Client {
	client := &Client{
		timeout: 30 * time.Second,
		providers: []provider.Provider{
			&provider.XdccEuProvider{},
			&provider.SunXdccProvider{},
			&provider.XdccRocksProvider{},
		},
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithProviders(providers []provider.Provider) Option {
	return func(c *Client) {
		c.providers = providers
	}
}

func (c *Client) Search(keywords []string) (search.Results, error) {
	results := search.Results{}

	if len(keywords) == 0 {
		return nil, errors.New("no keywords provided")
	}

	query := strings.Join(keywords, "+")
	resultsChannel := make(chan search.Results)
	defer close(resultsChannel)

	ctx := context.Background()

	errChan := make(chan error)

	for _, p := range c.providers {
		ctx, cancel := context.WithTimeout(ctx, c.timeout)

		go func() {
			defer cancel()
			c.runSearch(p, resultsChannel, query)
			select {
			case <-ctx.Done():
				errChan <- errors.New(fmt.Sprintf("Provider %#v exceeded timeout\n", p))
			}
		}()
	}

	var err error

	for range c.providers {
		select {
		case r := <-resultsChannel:
			results = append(results, r...)
		case e := <-errChan:
			err = e
		}
	}

	return results, err
}

func (c *Client) runSearch(p provider.Provider, resultsChannel chan search.Results, query string) {
	resList, err := p.Search(query)
	if err != nil {
		resultsChannel <- search.Results{}
		return
	}

	resultsChannel <- resList
}

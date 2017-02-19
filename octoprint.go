// Package octoprint facilitates communcation to an octoprint server.
package octoprint

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/dustin/httputil"
)

// A Client to an octoprint server.
type Client struct {
	base  *url.URL
	token string
}

func (c *Client) url(path string) *url.URL {
	u := *c.base
	u.Path = path
	return &u
}

func (c *Client) do(ctx context.Context, method, path string, r io.ReadCloser) (io.ReadCloser, error) {
	u := c.url(path)
	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("X-Api-Key", c.token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		return nil, httputil.HTTPError(res)
	}
	return res.Body, nil

}

func (c *Client) fetch(ctx context.Context, path string) (io.ReadCloser, error) {
	return c.do(ctx, "GET", path, nil)
}

func (c *Client) fetchJSON(ctx context.Context, path string, o interface{}) error {
	r, err := c.fetch(ctx, path)
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(o)
}

// New creates a connection to an octoprint server.
func New(base, token string) (*Client, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return &Client{u, token}, nil
}

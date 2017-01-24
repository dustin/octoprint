// Package octoprint facilitates communcation to an octoprint server.
package octoprint

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

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

// TimelapseConfig represents the configuration of tapelapse recording.
type TimelapseConfig struct {
	CapturePostRoll bool   `json:"capturePostRoll"`
	FPS             int    `json:"fps"`
	Interval        int    `json:"interval"`
	PostRoll        int    `json:"postRoll"`
	Type            string `json:"type"`
}

// A Timelapse entry contains all the fields representing a timelapse recording.
type Timelapse struct {
	Size    int64  `json:"bytes"`
	DateStr string `json:"date"`
	Name    string `json:"name"`
	SizeStr string `json:"size"`
	Path    string `json:"url"`

	c *Client
}

// URL returns the URL to the timelapse video on octoprint.
func (t Timelapse) URL() *url.URL {
	return t.c.url(t.Path)
}

// Fetch a timelapse video from octoprint.
func (t Timelapse) Fetch(ctx context.Context) (io.ReadCloser, error) {
	return t.c.fetch(ctx, t.Path)
}

// Delete a timelapse video from the octoprint server.
func (t Timelapse) Delete(ctx context.Context) error {
	r, err := t.c.do(ctx, "DELETE", "/api/timelapse/"+t.Name, nil)
	if err != nil {
		return err
	}
	return r.Close()
}

// ListTimelapses lists all of the available timelapse videos on the octoprint server.
func (c *Client) ListTimelapses(ctx context.Context) (*TimelapseConfig, []Timelapse, error) {
	v := struct {
		Config TimelapseConfig
		Files  []Timelapse
	}{}

	if err := c.fetchJSON(ctx, "/api/timelapse", &v); err != nil {
		return nil, nil, err
	}

	for i := range v.Files {
		v.Files[i].c = c
	}

	return &v.Config, v.Files, nil
}

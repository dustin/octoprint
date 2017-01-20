package octoprint

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/dustin/httputil"
)

type Client struct {
	base  *url.URL
	token string
}

func (c *Client) URL(path string) *url.URL {
	u := *c.base
	u.Path = path
	return &u
}

func (c *Client) fetch(path string) (io.ReadCloser, error) {
	u := c.URL(path)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		return nil, httputil.HTTPError(res)
	}
	return res.Body, nil
}

func (c *Client) fetchJSON(path string, o interface{}) error {
	r, err := c.fetch(path)
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(o)
}

func New(base, token string) (*Client, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return &Client{u, token}, nil
}

type TimelapseConfig struct {
	CapturePostRoll bool   `json:"capturePostRoll"`
	FPS             int    `json:"fps"`
	Interval        int    `json:"interval"`
	PostRoll        int    `json:"postRoll"`
	Type            string `json:"type"`
}

type Timelapse struct {
	Size    int64  `json:"bytes"`
	DateStr string `json:"date"`
	Name    string `json:"name"`
	SizeStr string `json:"size"`
	Path    string `json:"url"`

	c *Client
}

func (t Timelapse) URL() *url.URL {
	return t.c.URL(t.Path)
}

func (t Timelapse) Fetch() (io.ReadCloser, error) {
	return t.c.fetch(t.Path)
}

func (c *Client) ListTimelapses() (*TimelapseConfig, []Timelapse, error) {
	v := struct {
		Config TimelapseConfig
		Files  []Timelapse
	}{}

	if err := c.fetchJSON("/api/timelapse", &v); err != nil {
		return nil, nil, err
	}

	for i := range v.Files {
		v.Files[i].c = c
	}

	return &v.Config, v.Files, nil
}

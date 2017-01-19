package octoprint

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/dustin/httputil"
)

type Client struct {
	base  *url.URL
	token string
}

func (c *Client) fetch(path string, o interface{}) error {
	u := *c.base
	u.Path = path
	log.Printf("Fetching from %v %v", u.String(), c.token)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", c.token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return httputil.HTTPError(res)
	}

	d := json.NewDecoder(res.Body)
	return d.Decode(o)

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
	Bytes   int    `json:"bytes"`
	DateStr string `json:"date"`
	Name    string `json:"name"`
	SizeStr string `json:"size"`
	URL     string `json:"url"`
}

func (c *Client) ListTimelapses() (*TimelapseConfig, []Timelapse, error) {
	v := struct {
		Config TimelapseConfig
		Files  []Timelapse
	}{}

	if err := c.fetch("/api/timelapse", &v); err != nil {
		return nil, nil, err
	}

	return &v.Config, v.Files, nil
}

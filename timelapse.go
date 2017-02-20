package octoprint

import (
	"context"
	"io"
	"net/url"
)

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
	return t.c.url(t.Path, "")
}

// Fetch a timelapse video from octoprint.
func (t Timelapse) Fetch(ctx context.Context) (io.ReadCloser, error) {
	return t.c.fetch(ctx, t.Path, "")
}

// Delete a timelapse video from the octoprint server.
func (t Timelapse) Delete(ctx context.Context) error {
	r, err := t.c.do(ctx, "DELETE", "/api/timelapse/"+t.Name, "", nil)
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

	if err := c.fetchJSON(ctx, "/api/timelapse", "", &v); err != nil {
		return nil, nil, err
	}

	for i := range v.Files {
		v.Files[i].c = c
	}

	return &v.Config, v.Files, nil
}

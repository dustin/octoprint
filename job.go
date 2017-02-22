package octoprint

import "golang.org/x/net/context"

type JobState struct {
	Job struct {
		AveragePrintTime   *int `json:"averagePrintTime"`
		EstimatedPrintTime *int `json:"estimatedPrintTime"`
		Filament           *string
		File               struct {
			Timestamp int `json:"date"`
			Name      string
			Origin    string
			Path      string
			Size      int
		}
		LastPrintTimeStamp *int `json:"lastPrintTime"`
	}
	Progress struct {
		Completion          float64
		Filepos             int
		PrintTime           int
		PrintTimeLeft       int
		PrintTimeLeftOrigin string
	} `json:"progress"`
	State string
}

func (c *Client) JobState(ctx context.Context) (*JobState, error) {
	st := &JobState{}
	if err := c.fetchJSON(ctx, "/api/job", "", st); err != nil {
		return nil, err
	}
	return st, nil
}

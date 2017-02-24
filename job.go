package octoprint

import "golang.org/x/net/context"

type JobState struct {
	Job struct {
		AveragePrintTime   *float64 `json:"averagePrintTime"`
		EstimatedPrintTime *float64 `json:"estimatedPrintTime"`
		Filament           struct {
			Tool0 *struct {
				Length float64
				Volume float64
			}
			Tool1 *struct {
				Length float64
				Volume float64
			}
		}
		File struct {
			Timestamp int `json:"date"`
			Name      string
			Origin    string
			Path      string
			Size      int
		}
		LastPrintTimeStamp *float64 `json:"lastPrintTime"`
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

// JobState returns the current job state from octoprint.
func (c *Client) JobState(ctx context.Context) (*JobState, error) {
	st := &JobState{}
	if err := c.fetchJSON(ctx, "/api/job", "", st); err != nil {
		return nil, err
	}
	return st, nil
}

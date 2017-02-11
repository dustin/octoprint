package main

import (
	"context"

	"github.com/dustin/httputil"
	"github.com/dustin/octoprint/tool"
)

func main() {
	httputil.InitHTTPTracker(false)
	tool.ToolMain(context.Background(),
		map[string]tool.Command{
			"ls": {0, lsTimelineCmd, "", lsTimelineFlags},
			"dl": {1, dlTimelineCmd, "", dlTimelineFlags},
		})
}

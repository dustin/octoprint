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
			"state": {0, showState, "", showStateFlags},
		})
}

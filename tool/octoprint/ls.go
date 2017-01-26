package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dustin/octoprint"
	"golang.org/x/net/context"
)

var (
	lsTimelineFlags = flag.NewFlagSet("ls", flag.ExitOnError)
)

func lsTimelineCmd(ctx context.Context, c *octoprint.Client, args []string) {
	_, tls, err := c.ListTimelapses(ctx)
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	for _, tl := range tls {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\n", tl.Name, tl.DateStr, tl.SizeStr, tl.URL())
	}
	tw.Flush()
}

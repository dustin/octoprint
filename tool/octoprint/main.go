package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"github.com/dustin/httputil"
	"github.com/dustin/octoprint"
	"github.com/dustin/octoprint/tool"
)

var (
	lsTimelineFlags = flag.NewFlagSet("ls", flag.ExitOnError)
	dlTimelineFlags = flag.NewFlagSet("dl", flag.ExitOnError)

	dlConc = dlTimelineFlags.Int("concurrency", 4, "how many concurrent downloads to perform")
	dlRm   = dlTimelineFlags.Bool("rm", false, "delete already synced items")
)

func lsTimelineCmd(c *octoprint.Client, args []string) {
	_, tls, err := c.ListTimelapses()
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	for _, tl := range tls {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\n", tl.Name, tl.DateStr, tl.SizeStr, tl.URL())
	}
	tw.Flush()
}

func dlTimelineCmd(c *octoprint.Client, args []string) {
	_, tls, err := c.ListTimelapses()
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	grp, _ := errgroup.WithContext(context.Background())

	sem := make(chan bool, *dlConc)
	for _, tl := range tls {
		tl := tl
		grp.Go(func() error {
			sem <- true
			defer func() { <-sem }()
			dest := filepath.Join(args[0], tl.Name)

			st, err := os.Stat(dest)
			if err == nil && st.Size() == tl.Size {
				if *dlRm {
					log.Printf("Deleting (already present) %v", tl.Name)
					return tl.Delete()
				}
				return nil
			}

			log.Printf("Downloading %v -> %v (%v)", tl.Name, dest, tl.Size)

			f, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			defer f.Close()

			r, err := tl.Fetch()
			if err != nil {
				defer os.Remove(dest)
				return err
			}
			defer r.Close()
			_, err = io.Copy(f, r)
			if err != nil {
				return err
			}
			if *dlRm {
				log.Printf("Deleting %v", tl.Name)
				return tl.Delete()
			}
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		log.Fatalf("Error downloading: %v", err)
	}
}

func main() {
	httputil.InitHTTPTracker(false)
	tool.ToolMain(
		map[string]tool.Command{
			"ls": {0, lsTimelineCmd, "", lsTimelineFlags},
			"dl": {1, dlTimelineCmd, "", dlTimelineFlags},
		})
}

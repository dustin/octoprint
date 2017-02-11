package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dustin/octoprint"
	"golang.org/x/sync/errgroup"
)

var (
	dlTimelineFlags = flag.NewFlagSet("dl", flag.ExitOnError)

	dlConc = dlTimelineFlags.Int("concurrency", 4, "how many concurrent downloads to perform")
	dlRm   = dlTimelineFlags.Bool("rm", false, "delete already synced items")
)

func dlTimelineCmd(ctx context.Context, c *octoprint.Client, args []string) {
	_, tls, err := c.ListTimelapses(ctx)
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	grp, _ := errgroup.WithContext(ctx)

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
					return tl.Delete(ctx)
				}
				return nil
			}

			log.Printf("Downloading %v -> %v (%v)", tl.Name, dest, tl.Size)

			f, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			defer f.Close()

			r, err := tl.Fetch(ctx)
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
				return tl.Delete(ctx)
			}
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		log.Fatalf("Error downloading: %v", err)
	}
}

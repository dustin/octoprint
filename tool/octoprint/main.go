package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/dustin/octoprint"
)

var token = flag.String("token", "", "octoprint token")

type Command struct {
	Nargs  int
	F      func(c *octoprint.Client, args []string)
	Argstr string
	Flags  *flag.FlagSet
}

func (c Command) Usage(name string) {
	fmt.Fprintf(os.Stderr, "Usage:  %s %s\n", name, c.Argstr)
	if c.Flags != nil {
		os.Stderr.Write([]byte{'\n'})
		c.Flags.PrintDefaults()
	}
	os.Exit(64)
}

func setUsage(commands map[string]Command) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage:\n  %s [http://octopi/] cmd [-opts] cmdargs\n",
			os.Args[0])

		fmt.Fprintf(os.Stderr, "\nCommands:\n")

		ss := sort.StringSlice{}
		for k := range commands {
			ss = append(ss, k)
		}
		ss.Sort()

		for _, k := range ss {
			fmt.Fprintf(os.Stderr, "  %s %s\n", k, commands[k].Argstr)
		}

		fmt.Fprintf(os.Stderr, "\n---- Subcommand Options ----\n")

		for _, k := range ss {
			if commands[k].Flags != nil {
				fmt.Fprintf(os.Stderr, "\n%s:\n", k)
				commands[k].Flags.PrintDefaults()
			}
		}

		os.Exit(1)
	}
}

func MaybeFatal(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}

func Verbose(v bool, f string, a ...interface{}) {
	if v {
		log.Printf(f, a...)
	}
}

func ParseURL(ustr string) *url.URL {
	u, err := url.Parse(ustr)
	MaybeFatal(err, "Error parsing URL: %v", err)
	return u
}

func toolMain(commands map[string]Command) {
	log.SetFlags(log.Lmicroseconds)

	setUsage(commands)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
	}

	off := 0
	u := "http://octopi/"

	if strings.HasPrefix(flag.Arg(0), "http://") {
		u = flag.Arg(0)
		off++
	}

	c, err := octoprint.New(u, *token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating octoprint client: %v", err)
		os.Exit(1)
	}

	cmdName := flag.Arg(off)
	cmd, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %v\n", cmdName)
		flag.Usage()
	}

	args := flag.Args()[off+1:]
	nargs := len(args)
	if cmd.Flags != nil {
		cmd.Flags.Parse(args)
		nargs = cmd.Flags.NArg()
	}

	if cmd.Nargs == 0 {
	} else if cmd.Nargs < 0 {
		reqargs := -cmd.Nargs
		if nargs < reqargs {
			cmd.Usage(cmdName)
		}
	} else {
		if nargs != cmd.Nargs {
			cmd.Usage(cmdName)
		}
	}

	cmd.F(c, args)
}

var lsTimelineFlags = flag.NewFlagSet("ls", flag.ExitOnError)

func lsTimelineCmd(c *octoprint.Client, args []string) {
	_, tls, err := c.ListTimelapses()
	if err != nil {
		log.Fatalf("Error listing timelapses: %v", err)
	}
	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	for _, tl := range tls {
		fmt.Fprintf(tw, "%v\t%v\t%v\n", tl.Name, tl.DateStr, tl.SizeStr)
	}
	tw.Flush()
}

func main() {
	toolMain(
		map[string]Command{
			"ls": {0, lsTimelineCmd, "", lsTimelineFlags},
		})
}

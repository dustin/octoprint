// Package tool provides common functionality for octoprint commandline tools.
package tool

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/dustin/octoprint"
)

var token = flag.String("token", "", "octoprint token")

// Command represents a single command and its arguments.
type Command struct {
	Nargs  int
	F      func(context.Context, *octoprint.Client, []string)
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

// MaybeFatal terminates the program with the given error if the error is not nil.
func MaybeFatal(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}

// Verbose logs if v is true.
func Verbose(v bool, f string, a ...interface{}) {
	if v {
		log.Printf(f, a...)
	}
}

// ParseURL parses a URL with a fatal error if the URL can't be parsed.
func ParseURL(ustr string) *url.URL {
	u, err := url.Parse(ustr)
	MaybeFatal(err, "Error parsing URL: %v", err)
	return u
}

// ToolMain initializes all of the commandline definitions passed in,
// creates the client for octopi and dispatches to the appropriate
// commandline tool.
func ToolMain(ctx context.Context, commands map[string]Command) {
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

	cmd.F(ctx, c, cmd.Flags.Args())
}

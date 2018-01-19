package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/mqu/go-notify"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		simple  bool
		version bool
		help    bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&simple, "simple", false, "(shortcut: s) Simple output which doesn't show remained seconds.")
	flags.BoolVar(&simple, "s", false, "(shortcut: s) Simple output which doesn't show remained seconds.")
	flags.BoolVar(&version, "version", false, "(shortcut: v) Print version information and quit.")
	flags.BoolVar(&version, "v", false, "(shortcut: v) Print version information and quit.")
	flags.BoolVar(&help, "help", false, "(shortcut: h) Print this message.")
	flags.BoolVar(&help, "h", false, "(shortcut: h) Print this message.")
	flags.Usage = printUsage

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if help {
		printUsage()
		return ExitCodeOK
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", name, appVersion)
		return ExitCodeOK
	}

	d, err := getDuration(flags.Args())
	if err != nil {
		fmt.Fprintf(cli.errStream, "\033[31;1m%v\n", err)
		fmt.Fprintf(cli.errStream, "\033[31;1mPlease check usage (%s -h)\033[0m\n", name)
		return ExitCodeError
	}
	rem := int(d.Seconds())

	notify.Init("time-to-go")
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	defer close(sigCh)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	stop := make(chan bool)
	defer close(stop)
	fmt.Printf("Sleeping %v\n", strings.Replace(d.String(), "m", "min", 1))
	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				if simple != true {
					rem--
					if rem != 0 {
						min := rem / 60
						sec := rem % 60
						hour := min / 60
						min = min % 60
						timeLimit := ""
						switch {
						case hour > 0:
							timeLimit = fmt.Sprintf("%02vh%02vmin%02vs", hour, min, sec)
						case min > 0:
							timeLimit = fmt.Sprintf("   %02vmin%02vs", min, sec)
						default:
							timeLimit = fmt.Sprintf("        %02vs", sec)
						}
						fmt.Fprintf(cli.outStream, "\r%s remains...", timeLimit)
					} else {
						fmt.Fprintf(cli.outStream, "\r  0 sec(s) remains...\n")
						break loop
					}
				}
			case <-stop:
				break loop
			}
		}
	}()

	select {
	case <-sigCh:
		fmt.Fprintf(cli.errStream, "\nCancelled.\n")
		return ExitCodeOK
	case <-time.After(d):
		break
	}

	var g sync.WaitGroup
	g.Add(2)
	go func() {
		notify.Init("time-to-go")
		n := notify.NotificationNew("time-to-go", "Wake up!!!!", "appointment-soon")
		n.Show()
		g.Done()
	}()
	go func() {
		flashScreen(6)
		g.Done()
	}()

	g.Wait()
	return ExitCodeOK
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var helpMessage = `Usage:
  time-to-go <TIME>

Options:
  -s, --simple
        Simple output which doesn't show remained seconds.
  -h, --help
        Print this help message.
  -v, --version
        Print version information and quit.

Example:

To set alarm 3 minutes 20 seconds.

  $ time-to-go 3:20

time-to-go accepts TIME as the below format.
When you specify time unit, it must be one of units defined by International System of Units (SI) or units outside the SI.
i.e.) second: s, minute: min

  45 seconds: 45 s, 45s, .45, :45
  3 minutes: 3 min, 3min, 3.00, 3.0, 3. 3:00, 3:0, 3:
  2 minutes 40 seconds: 2 min 40 s, 2min 40s, 2 40, 2.40, 2:40

Press Ctrl+C to cancel the timer.
`

var re = regexp.MustCompile(`:+`)

// printUsage prints help message.
func printUsage() {
	fmt.Fprintf(os.Stderr, helpMessage)
	flag.PrintDefaults()
}

// getDuration converts args to time.Duration.
// It accepts various format. If conversion fails, corresponding error
// is returned.
func getDuration(args []string) (time.Duration, error) {
	var d time.Duration
	argString := strings.Join(args, ":")

	// 1. Replace "." with ":".
	argString = strings.Replace(argString, ".", ":", -1)
	// 2. Replace "min" with ":".
	argString = strings.Replace(argString, "min", ":", 1)
	// 3. Replace "::" with ":".
	argString = re.ReplaceAllString(argString, ":")
	// 4. Trim appended "s".
	if strings.HasSuffix(argString, ":s") {
		argString = strings.TrimRight(argString, ":s")
	}
	if strings.HasSuffix(argString, "s") {
		argString = strings.TrimRight(argString, "s")
	}
	// 5. Prepend ":" if it consists of numbers only.
	_, err := strconv.Atoi(argString)
	if err == nil {
		argString = ":" + argString
	}

	err = nil
	array := strings.Split(argString, ":")
	switch {
	case len(array) != 2:
		err = errors.New("Wrong format")
	case array[0] == "":
		s, err := strconv.Atoi(array[1])
		if err != nil {
			break
		}
		d = time.Duration(s) * time.Second
	case array[1] == "":
		m, err := strconv.Atoi(array[0])
		if err != nil {
			break
		}
		d = time.Duration(m) * time.Minute
	default:
		m, err := strconv.Atoi(array[0])
		if err != nil {
			break
		}
		s, err := strconv.Atoi(array[1])
		if err != nil {
			break
		}
		d = time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	}
	return d, err
}

// flashScreen makes current terminal screen flashing for t times
// using ANSI escape sequences. It has the potential not to work in
// some terminal environments.
func flashScreen(t int) {
	// workaround for TMUX.
	escSeq := "\x1b[?%s"
	_, tmux := os.LookupEnv("TMUX")
	if tmux == true {
		escSeq = "\x1bPtmux;\x1b\x1b[?%s\x1b\\"
	}
	for i := 0; i < t; i++ {
		// reverse video
		fmt.Fprintf(os.Stdout, escSeq, "5h")
		time.Sleep(500 * time.Millisecond)
		// normal video
		fmt.Fprintf(os.Stdout, escSeq, "5l")
		if i == t-1 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

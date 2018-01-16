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
i.e.) second: s, minute: min, hour: h

  45 seconds: 45 s, 45s, .45, :45
  3 minutes: 3 min, 3min, 3.00, 3.0, 3. 3:00, 3:0, 3:
  2 minutes 40 seconds: 2 min 40 s, 2min 40s, 2 40, 2.40, 2:40
  1 hours 15 minutes: 1 h 15 min, 1h 15min, 1.15.0, 1.15.00, 1:15:0, 1:15:00
  1 hours 20 minutes 30 seconds: 1 h 20 min 30 s, 1h 20min 30s, 1 20 30, 1.20.30, 1:20:30
  2 hours 40 seconds: 2 h 40 s, 2h 40s, 2 0 45

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
	// 2. Check the exisistence of "h", "min" and "s".
	if strings.Contains(argString, "h") && !strings.Contains(argString, "min") && strings.Contains(argString, "s") {
		argString = strings.Replace(argString, "h:", "h:0:", 1)
	}
	// 2. Replace "h" with ":".
	argString = strings.Replace(argString, "h", ":", 1)
	// 3. Replace "min" with ":".
	argString = strings.Replace(argString, "min", ":", 1)
	// 4. Replace "::" with ":".
	argString = re.ReplaceAllString(argString, ":")
	// 5. Trim appended "s".
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
	switch len(array) {
	case 2:
		switch {
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
	case 3:
		switch {
		case array[2] == "":
			h, err := strconv.Atoi(array[0])
			if err != nil {
				break
			}
			m, err := strconv.Atoi(array[1])
			if err != nil {
				break
			}
			d = time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
		default:
			h, err := strconv.Atoi(array[0])
			if err != nil {
				break
			}
			m, err := strconv.Atoi(array[1])
			if err != nil {
				break
			}
			s, err := strconv.Atoi(array[2])
			if err != nil {
				break
			}
			d = time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
		}
	default:
		err = errors.New("Wrong format")
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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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
  3 minutes: 3 min, 3min, 3.00, 3.0, 3. 3:00, 3:0, 3
  2 minutes 40 seconds: 2 min 40 s, 2min 40s, 2 40, 2.40, 2:40

Press Ctrl+C to cancel the timer.
`

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
	argsLen := len(args)
L1:
	switch {
	case argsLen == 0:
		return d, errors.New("Too few argument")
	case argsLen == 1:
		var err error
		// 1. Check if arg contains "." or ":".
		for _, sep := range []string{".", ":"} {
			array := strings.Split(args[0], sep)
			if len(array) != 2 {
				err = errors.New("Wrong format")
				continue
			}
			if array[0] == "" {
				s, err := strconv.Atoi(array[1])
				if err != nil {
					continue
				}
				d += time.Duration(s) * time.Second
				break L1
			} else if array[1] == "" {
				m, err := strconv.Atoi(array[0])
				if err != nil {
					continue
				}
				d += time.Duration(m) * time.Minute
				break L1
			} else {
				m, err := strconv.Atoi(array[0])
				if err != nil {
					continue
				}
				s, err := strconv.Atoi(array[1])
				if err != nil {
					continue
				}
				d += time.Duration(m)*time.Minute + time.Duration(s)*time.Second
				break L1
			}
		}
		// 2. Check if arg ends with "s".
		lastChar := args[0][len(args[0])-1]
		if lastChar == 's' {
			s, err := strconv.Atoi(args[0][:len(args[0])-1])
			if err == nil {
				d += time.Duration(s) * time.Second
				break L1
			}
		}
		// 3. Check if arg ends with "min".
		if len(args[0]) >= 4 {
			last3Chars := args[0][len(args[0])-3:]
			if last3Chars == "min" {
				m, err := strconv.Atoi(args[0][:len(args[0])-3])
				if err == nil {
					d += time.Duration(m) * time.Minute
					break L1
				}
			}
		}
		if err != nil {
			return d, err
		}
	case argsLen == 2:
		// 1. Check if 1st arg is number.
		arg0, err := strconv.Atoi(args[0])
		if err == nil {
			switch args[1] {
			case "s":
				d += time.Duration(arg0) * time.Second
				break L1
			case "min":
				d += time.Duration(arg0) * time.Minute
				break L1
			default:
				arg1, err := strconv.Atoi(args[1])
				if err != nil {
					return d, err
				}
				d += time.Duration(arg0) * time.Minute
				d += time.Duration(arg1) * time.Second
				break L1
			}
		}
		// 2. Check if the args are like "XXmin YYs".
		if len(args[0]) < 4 {
			return d, errors.New("Wrong format")
		}
		arg0last3Chars := args[0][len(args[0])-3:]
		if arg0last3Chars != "min" {
			return d, errors.New("Wrong format")
		}
		arg1lastChar := args[1][len(args[1])-1]
		if arg1lastChar != 's' {
			return d, errors.New("Wrong format")
		}
		m, err := strconv.Atoi(args[0][:len(args[0])-3])
		if err != nil {
			return d, err
		}
		s, err := strconv.Atoi(args[1][:len(args[1])-1])
		if err != nil {
			return d, err
		}
		d += time.Duration(m) * time.Minute
		d += time.Duration(s) * time.Second
	case argsLen == 4:
		if args[1] != "min" {
			return d, errors.New("Wrong format")
		}
		if argsLen == 4 && args[3] != "s" {
			return d, errors.New("Wrong format")
		}
		// At least, 1st arg must be number.
		m, err := strconv.Atoi(args[0])
		if err != nil {
			return d, err
		}
		// 4th arg also must be number.
		s, err := strconv.Atoi(args[2])
		if err != nil {
			return d, err
		}
		d += time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	case argsLen == 3 || argsLen > 4:
		return d, errors.New("Too much arguments")
	}
	return d, nil
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

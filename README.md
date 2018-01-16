# time-to-go

[![Build Status](https://travis-ci.org/tdsh/time-to-go.svg?branch=master)](https://travis-ci.org/tdsh/time-to-go)

Simple command to trigger an alarm (desktop notification and screen flashing) in specified duration.

## Description

time-to-go is a simple command. It accepts a duration to the alarm and when the duration passes, triggers an alarm with desktop notification and screen flashing. It can be used for any miscellaneous purpose, from meditation timer to cooking instant noodle.

## Usage

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

time-to-go accepts TIME as the below format. When you specify time unit, it must be one of units defined by International System of Units (SI) or units outside the SI. i.e.) second: s, minute: min

  45 seconds: 45 s, 45s, .45, :45
  3 minutes: 3 min, 3min, 3.00, 3.0, 3. 3:00, 3:0, 3
  2 minutes 40 seconds: 2 min 40 s, 2min 40s, 2 40, 2.40, 2:40
  1 hours 15 minutes: 1 h 15 min, 1h 15min, 1.15.0, 1.15.00, 1:15:0, 1:15:00
  1 hours 20 minutes 30 seconds: 1 h 20 min 30 s, 1h 20min 30s, 1 20 30, 1.20.30, 1:20:30
  2 hours 40 seconds: 2 h 40 s, 2h 40s, 2 0 45

Press Ctrl+C to cancel the timer.

## Install

time-to-go uses Libnotify. So libnotify and Go bindings [go-notify](https://github.com/mqu/go-notify) are required.

### libnotify

RedHat (Please replace by dnf in Fedora)

```bash
$ yum install libnotify-devel
```

Debian/Ubuntu

```bash
$ sudo apt-get install libnotify-dev
```

To install, use `go get`:

```bash
$ go get -d github.com/tdsh/time-to-go
```

## Contribution

1. Fork ([https://github.com/tdsh/time-to-go/fork](https://github.com/tdsh/time-to-go/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[tdsh](https://github.com/tdsh)

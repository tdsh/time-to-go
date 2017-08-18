# time-to-go

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

time-to-go accepts TIME as the below format.
  45 sec: 45 sec, 45 s, 45sec, 45s, .45, :45
    3 min: 3 min, 3 m, 3min, 3m, 3.00, 3.0, 3. 3:00, 3:0, 3
      2 min 40 sec: 2 min 40 sec, 2 m 40 s, 2 min 40, 2 m 40, 2 40, 2.40, 2:40

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

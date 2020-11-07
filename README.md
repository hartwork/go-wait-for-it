# go-wait-for-it

## About

**go-wait-for-it** is a mostly complete re-write of
Python [wait-for-it](https://github.com/clarketm/wait-for-it)
in Go 1.14+: as an *exercise* in writing code in Go and for *fun*.
To be clear, there is nothing wrong with the Python version: it's fine.

**go-wait-for-it** is my first full program written in Go
outside of [A Tour of Go](https://tour.golang.org/).
As such, you may find *beginner mistakes* in here, at least for the moment.
The goal is to have production level Go code in here, *eventually*.


## Build and run

```console
$ go version
go version go1.14.6 linux/amd64
$ go build .
$ ./go-wait-for-it --version
wait-for-it 1.0.0
```


## Examples

```console
$ ./go-wait-for-it -t 2 -s :631 -s localhost:631 -s 127.0.0.1:631 -- echo 'CUPS is available'
[*] Trying to connect to :631...
[*] Trying to connect to localhost:631...
[*] Trying to connect to 127.0.0.1:631...
[+] Connected to 127.0.0.1:631 after 201.653µs.
[+] Connected to :631 after 158.548µs.
[+] Connected to localhost:631 after 381.536µs.
[*] Running command: echo CUPS is available
CUPS is available
[+] Command succeeded.
```


## Usage

```console
$ go-wait-for-it --help
Wait for service(s) to be available before executing a command.

Usage:
  wait-for-it [flags] [-s|--service [HOST]:PORT]... [--] [COMMAND [ARG ..]]

Flags:
  -h, --help              help for wait-for-it
  -q, --quiet             do not output any status messages
  -s, --service strings   services to test (format '[HOST]:PORT')
  -t, --timeout uint      timeout in seconds, 0 for no timeout (default 15)
  -v, --version           version for wait-for-it
```


## Alternatives

### Go

- [`github.com/alioygur/wait-for`](https://github.com/alioygur/wait-for)
- [`github.com/maxcnunes/waitforit`](https://github.com/maxcnunes/waitforit)


### Python

- [`github.com/clarketm/wait-for-it`](https://github.com/clarketm/wait-for-it)
  — package `wait-for-it` [on PyPI](https://pypi.org/project/wait-for-it/)
- [`github.com/David-Lor/python-wait4it`](https://github.com/David-Lor/python-wait4it/)
  — package `wait4it` [on PyPI](https://pypi.org/project/wait4it/)


### Rust

- [`github.com/shenek/wait-for-them`](https://github.com/shenek/wait-for-them)
  — package `wait-for-them` [on crates.io](https://crates.io/crates/wait-for-them)


### Shell

- [`github.com/eficode/wait-for`](https://github.com/eficode/wait-for)
  — POSIX shell
- [`github.com/vishnubob/wait-for-it`](https://github.com/vishnubob/wait-for-it)
  — Bash, package `wait-for-it` [in Debian](https://packages.debian.org/unstable/wait-for-it)

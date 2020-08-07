# SANUS

[![build](https://github.com/vitrun/sanus/workflows/build/badge.svg)](https://github.com/vitrun/sanus/actions/)
[![Go Report Card](https://goreportcard.com/badge/github.com/vitrun/sanus)](https://goreportcard.com/report/github.com/vitrun/sanus)

A collection of Go Vet-style linters to make program sanus ([ˈsaː.nus], sound in body, healthy; sound in mind, sane).
- panicable. Check for possible panics that would result in program exit.

## Panicable

`Panicable` finds goroutines that have no `defer` or `recover` statements. Once they panic, not necessarily though, the whole program is endangered.


These patterns while be identified:

**pattern 1**
```go
func main() {
	go func() {  // no defer
		panic("I'm paniced")
	}()
}
```

**pattern 2**
```go
func a() {
	go func() { // no recover
		defer func() {
		}()
    panic("doomed")
	}()
}
```

Methods are supported:
```go
type S struct{}

func (s *S) foo() {
	go func() { // no defer
		println("foo")
	}()
}
```

## Install

Specify the cmd to install, `panicable` for example:

```
go get github.com/vitrun/sanus/cmd/panicable
```

This will install `panicable` to `$GOPATH/bin`, so make sure that it is included in your `$PATH` environment variable.


## Usage

Run cmd on current package like this:

```
$ panicable .
```

Or supply multiple packages, separated by spaces:

```
$ panicable example/cmd example/util strings
```

To check a package and, recursively, all its imports, use `./...`:

```
$ panicable net/...
```

all cmds in sanus accept the same flags as `go vet`:

```
Flags:
  -V    print version and exit
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -trace string
        write trace log to this file
```

## Development

To get the source code and compile/test specific cmd, run this:

```
$ git clone https://github.com/vitrun/sanus
$ cd sanus/panicable
$ go build -o panicable ./cmd/panicable
$ go test -v
```

`sanus` uses the testing infrastructure from `golang.org/x/tools/go/analysis/analysistest`. To add a test case, add new package in related `testdata/src`. 

In the test case source files, add annotation comments to the lines that should be reported (or not). The comments must
look like this:

```
func main() {
	go func() {  // want "no defer"
		panic("I'm paniced")
	}()
}
```

Annotations that indicate a line that should be reported must begin with a `want` followed by the desired message .

Since `sanus` is built upon the Go Vet standard infrastructure, you can import the passes into your own Go Vet-based linter.

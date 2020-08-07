// The panic command applies the github.com/vitrun/sanus/panicable
// analysis to the specified packages of Go source code.
package main

import (
	"github.com/vitrun/sanus/panicable"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(panicable.Analyzer) }

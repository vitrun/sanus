package panicable

import (
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = `check go routines without defer and recover statements for possible panics
This check reports conditions such as:

func main() {
	go func() {    // no defer 
		panic("I'm paniced")
	}()
}

func a() {
	go func() {
		defer func() {  // no recover
		}()
	}()
}
`

// Analyzer is ready to be used in a single checker
var Analyzer = &analysis.Analyzer{
	Name:     "panicable",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	ssainput := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	runPkg(pass, ssainput.Pkg)
	return nil, nil
}

func runPkg(pass *analysis.Pass, pkg *ssa.Package) {
	if pkg == nil {
		return
	}
	for _, m := range pkg.Members {
		switch mem := m.(type) {
		case *ssa.Function:
			detectByFunc(pass, mem)
		case *ssa.Type:
			mset := pkg.Prog.MethodSets.MethodSet(types.NewPointer(mem.Type()))
			for i, n := 0, mset.Len(); i < n; i++ {
				fun := pkg.Prog.MethodValue(mset.At(i))
				// point receiver
				if len(fun.Synthetic) == 0 {
					detectByFunc(pass, fun)
					continue
				}
				// value receiver
				for _, b := range fun.Blocks {
					// fixed instructions: Call, UnOp, Call and Return
					if len(b.Instrs) != 4 {
						continue
					}
					c, ok := b.Instrs[2].(*ssa.Call)
					if !ok {
						continue
					}
					if f, ok := c.Call.Value.(*ssa.Function); ok {
						detectByFunc(pass, f)
					}
				}
			}
		}
	}
}

func detectByFunc(pass *analysis.Pass, fun *ssa.Function) {
	for _, b := range fun.Blocks {
		detectByBlock(pass, b)
	}
}

func detectByBlock(pass *analysis.Pass, block *ssa.BasicBlock) {
	reportf := func(category string, pos token.Pos, msg string) {
		pass.Report(analysis.Diagnostic{
			Pos:      pos,
			Category: category,
			Message:  msg,
		})
	}

	for _, in := range block.Instrs {
		g, ok := in.(*ssa.Go)
		if !ok {
			continue
		}
		f, ok := g.Call.Value.(*ssa.Function)
		if !ok {
			continue
		}
		if f.Recover == nil {
			reportf("nodefer", f.Pos(), "no defer")
			continue
		}
		for _, af := range f.AnonFuncs {
			if hasRecover(af.Blocks[0].Instrs) {
				continue
			}
		}
		reportf("norecover", f.Pos(), "no recover")
	}
}

// hasRecover checks whether `recover()` exists in given instructions
func hasRecover(instrs []ssa.Instruction) bool {
	for _, inst := range instrs {
		if c, ok := inst.(*ssa.Call); ok {
			if v, ok := c.Call.Value.(*ssa.Builtin); ok && v.Name() == "recover" {
				return true
			}
		}
	}
	return false
}

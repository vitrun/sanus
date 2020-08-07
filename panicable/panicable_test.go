package panicable_test

import (
	"github.com/vitrun/sanus/panicable"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{
			name:     "no defer",
			patterns: []string{"a"},
		},
		{
			name:     "no recover",
			patterns: []string{"b"},
		},
		{
			name:     "named func",
			patterns: []string{"c"},
		},
		{
			name:     "func in another file",
			patterns: []string{"d"},
		},
		{
			name:     "in method",
			patterns: []string{"e"},
		},
	}
	testdata := analysistest.TestData()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysistest.Run(t, testdata, panicable.Analyzer, tt.patterns...)
		})
	}
}

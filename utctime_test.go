package utctime

import (
	"go/ast"
	"go/parser"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFromTestdata(t *testing.T) {
	testdata := analysistest.TestData()

	plugin, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("failed to build analyzers: %v", err)
	}

	analysistest.Run(t, testdata, analyzers[0], "utctime")
}

func TestIsTimeNowUTC(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "valid time.Now().UTC()",
			code:     "time.Now().UTC()",
			expected: true,
		},
		{
			name:     "wrong selector name",
			code:     "time.Now().Local()",
			expected: false,
		},
		{
			name:     "wrong package name",
			code:     "other.Now().UTC()",
			expected: false,
		},
		{
			name:     "wrong method name",
			code:     "time.Since().UTC()",
			expected: false,
		},
		{
			name:     "non-call expression before UTC",
			code:     "time.Now.UTC()",
			expected: false,
		},
		{
			name:     "non-selector expression in call",
			code:     "time().UTC()",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// The expression is a CallExpr with a SelectorExpr
			call, ok := node.(*ast.CallExpr)
			if !ok {
				t.Fatal("Expected CallExpr")
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				t.Fatal("Expected SelectorExpr")
			}

			got := isTimeNowUTC(sel)
			if got != tt.expected {
				t.Errorf("isTimeNowUTC() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package utctime

import (
	"go/ast"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const linterName = "utctime"
const linterDoc = "Checks that time.Now() is followed by .UTC()"

type Plugin struct{}

func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{&analysis.Analyzer{
		Name:     linterName,
		Doc:      linterDoc,
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func init() {
	register.Plugin(linterName, New)
}

// isTimeNowUTC checks if a selector expression represents time.Now().UTC().
func isTimeNowUTC(sel *ast.SelectorExpr) bool {
	if sel.Sel.Name != "UTC" {
		return false
	}

	nowCall, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	nowSel, ok := nowCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := nowSel.X.(*ast.Ident)
	if !ok || ident.Name != "time" {
		return false
	}

	return nowSel.Sel.Name == "Now"
}

// isTimeNow checks if a call expression represents time.Now().
func isTimeNow(call *ast.CallExpr) bool {
	selectorExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok || ident.Name != "time" {
		return false
	}

	return selectorExpr.Sel.Name == "Now"
}

// findParentNode checks if the current node is part of a larger expression.
func findParentNode(node ast.Node, file *ast.File) (ast.Node, bool) {
	var parent ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		if x, ok := n.(*ast.SelectorExpr); ok && x.X == node {
			parent = n
			return false
		}
		return true
	})
	return parent, parent != nil
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// First pass: find all time.Now().UTC() calls
		utcCalls := make(map[ast.Node]bool)
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if isTimeNowUTC(sel) {
				utcCalls[sel.X] = true
			}
			return true
		})

		// Second pass: find bare time.Now() calls
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !isTimeNow(call) {
				return true
			}

			// Skip if this Now() call is part of a UTC() call
			if parent, found := findParentNode(n, file); found {
				if parentSel, isSelector := parent.(*ast.SelectorExpr); isSelector && parentSel.Sel.Name == "UTC" {
					return true
				}
			}

			pass.Reportf(call.Pos(), "time.Now() must be followed by .UTC()")
			return true
		})
	}

	//nolint:nilnil // ignore for analysis testing.
	return nil, nil
}

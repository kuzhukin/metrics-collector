package noexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noexitcheck",
	Doc:  "check os.Exit callings in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				handleCallExpr(pass, x)
			}

			return true
		})
	}

	return nil, nil
}

func handleCallExpr(pass *analysis.Pass, fn *ast.FuncDecl) {
	if fn.Name.Name != "main" {
		return
	}

	v := &visitor{pass: pass}
	for _, stmt := range fn.Body.List {
		ast.Walk(v, stmt)
	}
}

type visitor struct {
	pass *analysis.Pass
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	if selectorExpr.Sel.Name == "Exit" {
		v.pass.Reportf(node.Pos(), "using os.Exit in main func is forbidden")
		return nil
	}

	return v
}

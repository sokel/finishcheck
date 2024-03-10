package pkg

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "finishcheck",
		Doc:      "Checks that whatever is opened and closed correctly",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

var funcStartSpanFromContext = "StartSpanFromContext"
var funcStartSpan = "StartSpan"
var funcFinish = "Finish"

func run(pass *analysis.Pass) (interface{}, error) {
	spanOpened := false
	var spanOpenStmt *ast.AssignStmt

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.AssignStmt:
				if len(stmt.Rhs) == 1 {
					if callExpr, ok := stmt.Rhs[0].(*ast.CallExpr); ok {
						if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if selExpr.Sel.Name == funcStartSpanFromContext ||
								selExpr.Sel.Name == funcStartSpan {
								if spanOpened {
									pass.Reportf(spanOpenStmt.Pos(), "span is opened but not closed")
								}

								spanOpened = true
								spanOpenStmt = stmt
							}
						}
					}
				}
			case *ast.ExprStmt:
				if callExpr, ok := stmt.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if selExpr.Sel.Name == funcFinish {
							if !spanOpened {
								pass.Reportf(stmt.Pos(), "span is finished before being opened")
							}

							spanOpened = false
						}
					}
				}
			}
			return true
		})
	}

	if spanOpened {
		pass.Reportf(spanOpenStmt.Pos(), "span is opened but not closed")
	}

	return nil, nil
}

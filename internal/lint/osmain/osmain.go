package osmain

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
)

var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "errexit",
	Doc:  "check for that user does'nt use os.exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if c, ok := n.(*ast.FuncDecl); ok {
				if c.Name.String() == "main" {
					log.Printf("%v", c.Name)
					ast.Inspect(c, func(n1 ast.Node) bool {
						if i, ok := n1.(*ast.CallExpr); ok {
							if m, ok := i.Fun.(*ast.SelectorExpr); ok {
								if m.Sel.Name == "Exit" {
									pass.Reportf(i.Pos(), "os.Exit used inside main()")
									return false
								}
							}
						}
						return true
					})
				}
			}

			return true
		})
	}
	return nil, nil
}

package slerrors

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func RunErrOSExit(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {
			// проверяем, какой конкретный тип лежит в узле
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					ast.Inspect(x, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
								if ident, ok := sel.X.(*ast.Ident); ok {
									if ident.Name == "os" && sel.Sel.Name == "Exit" {
										pos := pass.Fset.Position(call.Pos())
										fmt.Printf("  Found os.Exit call at %s\n", pos)
									}
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

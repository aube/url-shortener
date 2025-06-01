package slerrors

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RunErrOSExit analyzes Go AST to detect direct calls to os.Exit within main functions.
//
// This function traverses the AST of each file in the provided analysis pass,
// specifically looking for function declarations named "main". Within these main
// functions, it checks for any calls to os.Exit and reports their positions.
//
// Returns:
//   - nil, nil: This analyzer doesn't return any meaningful data, only detects patterns.
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
										pass.Reportf(ident.NamePos, "found os.Exit call")

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

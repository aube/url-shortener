package slerrors

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// errorType represents the Go error interface type for type comparison purposes.
var errorType = types.
	// ищем тип error в области вилимости Universe, в котором находятся
	// все предварительно объявленные объекты Go
	Universe.Lookup("error").
	// получаем объект, представляющий тип error
	Type().
	// получаем тип, при помощи которого определен тип error (см. https://go.dev/ref/spec#Underlying_types);
	// мы знаем, что error определен как интерфейс, приведем полученный объект к этому типу
	Underlying().(*types.Interface)

// isErrorType checks if a given type implements the error interface.
//
// Parameters:
//   - t: The type to check against the error interface
//
// Returns:
//   - bool: true if the type implements error, false otherwise
func isErrorType(t types.Type) bool {
	// проверяем, что t реализует интерфейс, при помощи которого определен тип error,
	// т.е. для типа t определен метод Error() string
	return types.Implements(t, errorType)
}

// resultErrors analyzes a function call expression to determine which return values are errors.
//
// Parameters:
//   - pass: The analysis pass containing type information
//   - call: The function call expression to analyze
//
// Returns:
//   - []bool: A slice where each element indicates whether the corresponding
//     return value is of error type
func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named: // возвращается значение
		return []bool{isErrorType(t)}
	case *types.Pointer: // возвращается указатель
		return []bool{isErrorType(t)}
	case *types.Tuple: // возвращается несколько значений
		s := make([]bool, t.Len())
		for i := range t.Len() {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

// isReturnError checks if a function call returns at least one error value.
//
// Parameters:
//   - pass: The analysis pass containing type information
//   - call: The function call expression to check
//
// Returns:
//   - bool: true if the call returns any error values, false otherwise
func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

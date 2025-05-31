package slerrors

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var errorType = types.
	// ищем тип error в области вилимости Universe, в котором находятся
	// все предварительно объявленные объекты Go
	Universe.Lookup("error").
	// получаем объект, представляющий тип error
	Type().
	// получаем тип, при помощи которого определен тип error (см. https://go.dev/ref/spec#Underlying_types);
	// мы знаем, что error определен как интерфейс, приведем полученный объект к этому типу
	Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	// проверяем, что t реализует интерфейс, при помощи которого определен тип error,
	// т.е. для типа t определен метод Error() string
	return types.Implements(t, errorType)
}

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

// isReturnError возвращает true, если среди возвращаемых значений есть ошибка.
func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

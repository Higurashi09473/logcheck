package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"github.com/Higurashi09473/logcheck/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for compliance with rules",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name

				if utils.IsLogMethod(methodName) {
					checkStartRune(pass, call)
					checkEnglishLanguage(pass, call)
					checkSpecialChars(pass, call)
					checkSensitiveData(pass, call)
				}
			}

			return true
		})
	}
	return nil, nil
}

func checkStartRune(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]

	lit, ok := firstArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	message, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	if message == "" {
		return
	}

	firstRune := []rune(message)[0]

	if unicode.IsUpper(firstRune) {
		pass.Reportf(lit.Pos(), "log message should not start with uppercase letter: %q", message)
	}
}

func checkEnglishLanguage(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]
	lit, ok := firstArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	// Правильное извлечение строки
	message, err := strconv.Unquote(lit.Value)
	if err != nil {
		return // Некорректная литерала — пропускаем (парсер уже ругнётся)
	}

	if message == "" {
		return
	}

	// Используем хелпер
	if !isEnglish(message) {
		pass.Reportf(lit.Pos(), "log message must be in English only (ASCII letters): %q", message)
	}
}


func isEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !utils.IsAsciiLatinLetter(r) {
			return false
		}
	}
	return true
}

func checkSpecialChars(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]
	lit, ok := firstArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	message, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	if message == "" {
		return
	}

	var offending []rune 
	for _, r := range message {
		if utils.IsEmoji(r) {
			offending = append(offending, r)
		} else if utils.IsForbiddenPunctuation(r) {
			offending = append(offending, r)
		}
	}

	if len(offending) > 0 {
		pass.Reportf(lit.Pos(), "log message must not contain special characters (found: %q): %q", string(offending), message)
	}
}

func checkSensitiveData(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}

		message, err := strconv.Unquote(lit.Value)
		if err != nil {
			continue
		}

		message = strings.TrimSpace(message) // Убираем пробелы для чистоты

		if message == "" {
			continue
		}

		if utils.ContainsSensitiveData(message) {
			pass.Reportf(lit.Pos(), "potential sensitive data in log message (e.g., passwords, tokens): %q", message)
		}
	}
}
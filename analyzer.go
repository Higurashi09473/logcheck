package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"github.com/Higurashi09473/logcheck/config"
	"github.com/Higurashi09473/logcheck/utils"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks log messages for compliance with rules",
	URL:      "https://github.com/Higurashi09473/logcheck",
	Flags:    config.NewFlagSet(),
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	flagMap := config.FlagsToMap(&pass.Analyzer.Flags)

	var cfg config.Options
	mapstructure.Decode(flagMap, &cfg)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name

				if utils.IsLogMethod(methodName) {
					if cfg.Lowercase {
						checkStartRune(pass, call)
					}
					if cfg.EnglishOnly {
						checkEnglishLanguage(pass, call)
					}
					if cfg.NoSpecialChars {
						checkSpecialChars(pass, call)
					}
					if cfg.NoSensitiveData {
						checkSensitiveData(pass, call)
					}
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

	if !unicode.IsUpper(firstRune) {
		return
	}

	startPos := lit.Pos() + token.Pos(1)   // +1 для "
	firstRuneLen := len(string(firstRune)) // длина руны в байтах (для UTF-8)

	// Фикс: заменить первую руну на lowercase
	lowerRune := unicode.ToLower(firstRune)
	newFirst := []byte(string(lowerRune))

	pass.Report(analysis.Diagnostic{
		Pos:     lit.Pos(),
		End:     lit.End(),
		Message: "log message should not start with uppercase letter",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Make first letter lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     startPos,
						End:     startPos + token.Pos(firstRuneLen),
						NewText: newFirst,
					},
				},
			},
		},
	})
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
		pass.Reportf(lit.Pos(), "log message must be in English only (ASCII letters)")
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
		pass.Reportf(lit.Pos(), "log message must not contain special characters (found: %q)", string(offending))
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
			pass.Reportf(lit.Pos(), "potential sensitive data in log message (e.g., passwords, tokens)")
		}
	}
}

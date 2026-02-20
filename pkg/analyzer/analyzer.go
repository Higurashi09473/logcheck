package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for compliance with rules",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		method := sel.Sel.Name
		if method != "Debug" && method != "Info" && method != "Warn" && method != "Error"{
			return
		}

		obj := pass.TypesInfo.ObjectOf(sel.Sel)
		fn, ok := obj.(*types.Func)
		if !ok {
			return
		}

		sig := fn.Type().(*types.Signature)

		var msgIndex int
		var isLogging bool

		if sig.Recv() == nil {
			// Potential slog package-level function
			if fn.Pkg() != nil && fn.Pkg().Path() == "log/slog" && len(call.Args) >= 1 {
				msgIndex = 0
				isLogging = true
			}
		} else {
			// Potential zap method
			recvType := sig.Recv().Type()
			if ptr, ok := recvType.(*types.Pointer); ok {
				recvType = ptr.Elem()
			}
			if named, ok := recvType.(*types.Named); ok && named.Obj().Pkg() != nil {
				pkgPath := named.Obj().Pkg().Path()
				typeName := named.Obj().Name()
				if pkgPath == "go.uber.org/zap" && (typeName == "Logger" || typeName == "SugaredLogger") && len(call.Args) >= 1 {
					msgIndex = 0
					isLogging = true
				}
			}
		}

		if !isLogging {
			return
		}

		msgExpr := call.Args[msgIndex]

		// Check for sensitive data
		checkSensitive(pass, msgExpr, call.Pos())

		// Collect string literals and check language and special chars
		var lits []string
		collectLiterals(msgExpr, &lits)
		for _, s := range lits {
			if !isEnglish(s) {
				pass.Reportf(msgExpr.Pos(), "log message must be in English")
			}
			if !noSpecial(s) {
				pass.Reportf(msgExpr.Pos(), "log message must not contain special symbols or emojis")
			}
		}

		// Check starting with lowercase
		if startRune, ok := getStartRune(msgExpr); ok {
			if !unicode.IsLower(startRune) {
				pass.Reportf(msgExpr.Pos(), "log message must start with a lowercase letter")
			}
		}
	})

	return nil, nil
}

var sensitiveKeywords = []string{"password", "apikey", "api_key", "token", "secret"}

func checkSensitive(pass *analysis.Pass, expr ast.Expr, pos token.Pos) {
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			nameLower := strings.ToLower(ident.Name)
			for _, kw := range sensitiveKeywords {
				if strings.Contains(nameLower, kw) {
					pass.Reportf(pos, "log message contains potentially sensitive data")
					break
				}
			}
		}
		return true
	})
}

func collectLiterals(expr ast.Expr, lits *[]string) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err == nil {
				*lits = append(*lits, s)
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			collectLiterals(e.X, lits)
			collectLiterals(e.Y, lits)
		}
	case *ast.ParenExpr:
		collectLiterals(e.X, lits)
	}
}

func isEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}

func noSpecial(s string) bool {
	allowedPunct := ".,:;-/()[]{} "
	for _, r := range s {
		if unicode.IsGraphic(r) && !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune(allowedPunct, r) {
			return false
		}
	}
	// Check repeated punctuation
	repeated := []string{"!!", "??", "..", "::", ",,", "--", ";;"}
	for _, rep := range repeated {
		if strings.Contains(s, rep) {
			return false
		}
	}
	// Rough emoji detection
	emojiRanges := "\u2600\u26FF\u2700\u27BF\u1F300\u1F5FF\u1F900\u1F9FF\u1F000\u1FFFF"
	for _, r := range s {
		if strings.ContainsRune(emojiRanges, r) {
			return false
		}
	}

	return true
}

func getStartRune(expr ast.Expr) (rune, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err != nil {
				return 0, false
			}
			for _, r := range s {
				if !unicode.IsSpace(r) {
					return r, true
				}
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return getStartRune(e.X)
		}
	case *ast.ParenExpr:
		return getStartRune(e.X)
	}
	return 0, false
}
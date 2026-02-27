package analyzer

import (
    "golang.org/x/tools/go/analysis"
)

func AnalyzerPlugin() []*analysis.Analyzer {
    return []*analysis.Analyzer{Analyzer}
}
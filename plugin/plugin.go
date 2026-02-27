package analyzer

import (
    "github.com/Higurashi09473/logcheck"
    "golang.org/x/tools/go/analysis"
)

func AnalyzerPlugin() []*analysis.Analyzer {
    return []*analysis.Analyzer{analyzer.Analyzer}
}
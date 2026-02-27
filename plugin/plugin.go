package main

import (
	"github.com/higurashi09473/logcheck/pkg/analyzer"
    "golang.org/x/tools/go/analysis"
)

func AnalyzerPlugin() []*analysis.Analyzer {
    return []*analysis.Analyzer{analyzer.Analyzer}
}
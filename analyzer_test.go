package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLowercase(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "lowercase")
}

func TestEnglish(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "english")
}

func TestSpecialChars(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "specialchars")
}

func TestSensitiveData(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "sensitive")
}

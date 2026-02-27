package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var logMethods = map[string]struct{}{
	// std log
	"Print":   {},
	"Println": {},
	"Printf":  {},

	// slog.Logger methods
	"Debug":        {},
	"DebugContext": {},
	"Info":         {},
	"InfoContext":  {},
	"Warn":         {},
	"WarnContext":  {},
	"Error":        {},
	"ErrorContext": {},
	"Log":          {},
	"LogAttrs":     {},

	// zap.Logger methods
	"DPanic": {},
	"Panic":  {},
	"Fatal":  {},

	// zap.SugaredLogger methods
	"Debugf":   {},
	"Infof":    {},
	"Warnf":    {},
	"Errorf":   {},
	"DPanicf":  {},
	"Panicf":   {},
	"Fatalf":   {},
	"Debugw":   {},
	"Infow":    {},
	"Warnw":    {},
	"Errorw":   {},
	"DPanicw":  {},
	"Panicw":   {},
	"Fatalw":   {},
	"Debugln":  {},
	"Infoln":   {},
	"Warnln":   {},
	"Errorln":  {},
	"DPanicln": {},
	"Panicln":  {},
	"Fatalln":  {},
}

func IsLogMethod(methodName string) bool {
	_, ok := logMethods[methodName]
	if ok {
		return true
	}

	if strings.ToLower(methodName) == "log" || strings.HasSuffix(strings.ToLower(methodName), "log") {
		return true
	}

	return false
}

func IsAsciiLatinLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func IsExtendedLatinLetter(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

func IsEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF)
}

func IsForbiddenPunctuation(r rune) bool {
	if !unicode.IsPunct(r) {
		return false
	}
	allowed := ".,?!:;\"'()-"
	return unicode.IsPunct(r) && !strings.ContainsRune(allowed, r)
}

var sensitiveRegex = regexp.MustCompile(`(?i)\b(password|passwd|pwd|token|jwt|bearer|api_key|apikey|api-key|secret|private_key|private-key|credit_card|card_number|cvv|ssn|social_security|authorization)\b`)

func ContainsSensitiveData(message string) bool {
	return sensitiveRegex.MatchString(message)
}
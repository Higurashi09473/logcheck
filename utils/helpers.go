package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var logMethods = map[string]struct{}{
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

	if strings.ToLower(methodName) == "log" {
		return true
	}

	return false
}

func IsAsciiLatinLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func IsEmoji(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F5FF) || // 1F300–1F5FF  — Miscellaneous Symbols and Pictographs
		(r >= 0x1F600 && r <= 0x1F64F) || // 1F600–1F64F  — Emoticons (основные смайлики)
		(r >= 0x1F680 && r <= 0x1F6FF) || // 1F680–1F6FF  — Transport and Map Symbols
		(r >= 0x1F900 && r <= 0x1F9FF) || // 1F900–1F9FF  — Supplemental Symbols and Pictographs
		(r >= 0x1FA70 && r <= 0x1FAFF) // 1FA70–1FAFF  — Symbols and Pictographs Extended-A + новые блоки
}

func IsForbiddenPunctuation(r rune) bool {
	if !unicode.IsPunct(r) {
		return false
	}
	allowed := ",:;\"'()- "
	return unicode.IsPunct(r) && !strings.ContainsRune(allowed, r)
}

var sensitiveRegex = regexp.MustCompile(`(?i)\b(password|passwd|pwd|token|jwt|bearer|api_key|apikey|api-key|secret|private_key|private-key|credit_card|card_number|cvv|ssn|social_security|authorization)\b`)

func ContainsSensitiveData(message string) bool {
	return sensitiveRegex.MatchString(message)
}

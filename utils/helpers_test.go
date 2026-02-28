package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLogMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"стандартный Print", "Print", false},
		{"Println", "Println", false},
		{"slog", "Info", true},
		{"slog context", "ErrorContext", true},
		{"zap", "DPanic", true},
		{"sugared zap", "Infow", true},
		{"sugared zap ln", "Fatalln", true},
		{"обычная функция", "Write", false},
		{"просто log", "log", true},
		{"не лог", "login", false},
		{"не лог 2", "logout", false},
		{"не лог 3", "dialog", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLogMethod(tt.input)
			assert.Equal(t, tt.expected, got, "input = %q", tt.input)
		})
	}
}

func TestIsAsciiLatinLetter(t *testing.T) {
	assert.True(t, IsAsciiLatinLetter('A'))
	assert.True(t, IsAsciiLatinLetter('z'))
	assert.True(t, IsAsciiLatinLetter('k'))
	assert.False(t, IsAsciiLatinLetter('1'))
	assert.False(t, IsAsciiLatinLetter(' '))
	assert.False(t, IsAsciiLatinLetter('é')) // не ascii
	assert.False(t, IsAsciiLatinLetter('ї')) // кириллица
	assert.False(t, IsAsciiLatinLetter('😀'))
}

func TestIsEmoji(t *testing.T) {
	assert.True(t, IsEmoji('😀'))
	assert.True(t, IsEmoji('🚀'))
	assert.True(t, IsEmoji('🧑')) // сложный эмодзи (комбинация)
	assert.True(t, IsEmoji('🥹'))
	assert.False(t, IsEmoji('A'))
	assert.False(t, IsEmoji('1'))
	assert.False(t, IsEmoji('-'))
	assert.False(t, IsEmoji('é'))
}

func TestIsForbiddenPunctuation(t *testing.T) {
	allowed := ".,?!:;\"'()- "
	for _, r := range allowed {
		assert.False(t, IsForbiddenPunctuation(r), "should be allowed: %q", r)
	}

	forbidden := []rune{'@', '#', '%', '&', '*', '/', '\\', '[', ']', '{', '}', '_'}
	for _, r := range forbidden {
		assert.True(t, IsForbiddenPunctuation(r), "should be forbidden: %q", r)
	}

	assert.False(t, IsForbiddenPunctuation('a')) // буква — не пунктуация
	assert.False(t, IsForbiddenPunctuation(' '))
}

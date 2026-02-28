package config

import (
	"flag"
)

type Options struct {
	Lowercase       bool `mapstructure:"check-start-rune"`
	EnglishOnly     bool `mapstructure:"check-english"`
	NoSpecialChars  bool `mapstructure:"check-special-chars"`
	NoSensitiveData bool `mapstructure:"check-sensitive"`
}

func NewFlagSet() flag.FlagSet {
	fs := flag.NewFlagSet("logcheck", flag.ContinueOnError)
	fs.Bool("check-start-rune", true, "check that log messages start with a lowercase letter")
	fs.Bool("check-english", true, "check that log messages contain only English characters")
	fs.Bool("check-special-chars", true, "check that log messages do not contain special characters or emoji")
	fs.Bool("check-sensitive", true, "check that log messages do not contain sensitive data")

	return *fs
}


func FlagsToMap(fs *flag.FlagSet) map[string]interface{} {
    m := make(map[string]interface{})
    fs.VisitAll(func(f *flag.Flag) {
        if b, ok := f.Value.(flag.Getter); ok {
            m[f.Name] = b.Get()
        } else {
            m[f.Name] = f.Value.String()
        }
    })
    return m
}
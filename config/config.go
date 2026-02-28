package config

import (
	"flag"
)

type Options struct {
	Lowercase       bool
	EnglishOnly     bool
	NoSpecialChars  bool
	NoSensitiveData bool
}


func NewFlagSet() flag.FlagSet {
	fs := flag.NewFlagSet("honeylog", flag.ContinueOnError)
	fs.Bool("lowercase", true, "check that log messages start with a lowercase letter")
	fs.Bool("english-only", true, "check that log messages contain only English characters")
	fs.Bool("no-special-chars", true, "check that log messages do not contain special characters or emoji")
	fs.Bool("no-sensitive-data", true, "check that log messages do not contain sensitive data")

	return *fs
}

func OptionsFromFlags(flags *flag.FlagSet) Options {
	return Options{
		Lowercase:         addFlag(flags, "lowercase"),
		EnglishOnly:       addFlag(flags, "english-only"),
		NoSpecialChars:    addFlag(flags, "no-special-chars"),
		NoSensitiveData:   addFlag(flags, "no-sensitive-data"),
	}
}

func addFlag(fs *flag.FlagSet, name string) bool {
	f := fs.Lookup(name)
	if f == nil {
		return true
	}
	if g, ok := f.Value.(flag.Getter); ok {
		if v, ok := g.Get().(bool); ok {
			return v
		}
	}
	return true
}
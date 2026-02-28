package sensitive

import "log/slog"

func example() {
	slog.Info("user password is secret123")                                 // want `potential sensitive data in log message`
	slog.Error("failed to validate apikey=sk0-live-abc123xyz")              // want `potential sensitive data in log message`
	slog.Warn("authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9") // want `potential sensitive data in log message`
	slog.Info("here is your jwt: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9")     // want `potential sensitive data in log message`
	slog.Debug("token received: abcdef123456")                              // want `potential sensitive data in log message`
	slog.Info("cvv is 456")                                                 // want `potential sensitive data in log message`
	slog.Warn("ssn 123-45-6789 found in payload")                           // want `potential sensitive data in log message`
	slog.Warn("authorization code 987654")                                  // want `potential sensitive data in log message`

	slog.Info("my passphrase is long enough")
	slog.Info("secretly planning something")
	slog.Info("tokenize payment")
	slog.Debug("pwd123 is not secure")
}

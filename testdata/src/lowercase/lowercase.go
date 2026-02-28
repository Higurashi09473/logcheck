package lowercase

import "log/slog"

func example() {
	slog.Info("server is starting on :8080")                    // OK
	slog.Info("Starting server on :8080")                       // want `log message should not start with uppercase letter`
	slog.Info("database connection failed after 3 retries")     // OK
	slog.Error("Failed to connect to database")                 // want `log message should not start with uppercase letter`
	slog.Warn("cannot open config file, using defaults")        // OK
	slog.Warn("Warning: low disk space")                        // want `log message should not start with uppercase letter`
	slog.Debug("received request from localhost")               // OK
	slog.Debug("Debug mode enabled")                            // want `log message should not start with uppercase letter`
	slog.Info("12345 requests processed successfully")          // OK (начинается с цифры)
	slog.Info("42")                                             // OK
	slog.Info("")                                               // OK (пустая строка)
	slog.Info("a single lowercase letter")                      // OK
}

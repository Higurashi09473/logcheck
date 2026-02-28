package specialchars

import "log/slog"

func example() {
	slog.Info("server started 🚀")             // want `log message must not contain special characters \(found: \"🚀\"\)`
	slog.Error("failed to connect 😞")         // want `log message must not contain special characters \(found: \"😞\"\)`
	slog.Info("all good 👍 great job")         // want `log message must not contain special characters \(found: \"👍\"\)`
	slog.Info("warning: something happened!") // want `log message must not contain special characters \(found: \"!\"\)`
	slog.Warn("please retry… later")          // want `log message must not contain special characters \(found: \"…\"\)`
	slog.Error("invalid data «abc»")          // want `log message must not contain special characters \(found: \"«»\"\)`
	slog.Info("price range: 100–200")         // want `log message must not contain special characters \(found: \"–\"\)`
	slog.Info("done! great job 🚀")            // want `log message must not contain special characters \(found: \"!🚀\"\)`
	slog.Error("failed 😢 — try again!")       // want `log message must not contain special characters \(found: \"😢—!\"\)`
	slog.Info("user: admin@example.com")      // want `log message must not contain special characters \(found: \"@.\"\)`
	slog.Info("version v2.3.1")               // want `log message must not contain special characters \(found: \"..\"\)`
	slog.Info("path: /api/v1/health")         // want `log message must not contain special characters \(found: \"///\"\)`
	slog.Info("status=200 text=\"ok\"")       // OK
}

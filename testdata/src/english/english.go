package english

import "log/slog"

func example() {
	slog.Info("стартуем сервер на 8080")         // want `log message must be in English only`
	slog.Info("starting server on port 8080")    // OK
	slog.Error("не удалось подключиться к базе") // want `log message must be in English only`
	slog.Debug("debugging connection")           // OK
	slog.Warn("ошибка авторизации пользователя") // want `log message must be in English only`
	slog.Info("123 numeric start")               // OK
	slog.Info("")                                // OK
	slog.Info("a")                               // OK
	slog.Info("а")                               // want `log message must be in English only`
}

package bad

import (
	"log/slog"

	"go.uber.org/zap"
)

func testSlog() {
	slog.Info("Starting server")  // want "log message must start with a lowercase letter"
	slog.Error("запуск сервера")  // want "log message must be in English"
	slog.Warn("server started!🚀") // want "log message must not contain special symbols or emojis"
	var password string
	slog.Info("user password: " + password) // want "log message contains potentially sensitive data"
	slog.Error("connection failed!!!")      // want "log message must not contain special symbols or emojis"
}

func testZap(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("Starting server") // want "log message must start with a lowercase letter"
	sugar.Error("ошибка")          // want "log message must be in English"
	logger.Warn("warning...")      // want "log message must not contain special symbols or emojis"
	var apiKey string
	logger.Info("api_key=" + apiKey) // want "log message contains potentially sensitive data"
}

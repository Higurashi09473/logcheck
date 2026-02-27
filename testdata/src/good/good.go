package good

import (
	"log/slog"

	"go.uber.org/zap"
)

func testSlog() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	slog.Info("user authenticated successfully")
}

func testZap(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("starting server")
	sugar.Error("failed to connect to database")
	logger.Warn("something went wrong")
	logger.Info("api request completed")
}

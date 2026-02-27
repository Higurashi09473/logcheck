package bad

import (
	"log/slog"

	"go.uber.org/zap"
)

func testSlog() {
	slog.Info("Starting server")  
	slog.Error("запуск сервера")  
	slog.Warn("server started!🚀") 
	var password string
	slog.Info("user password: " + password) 
	slog.Error("connection failed!!!")      
}

func testZap(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("Starting server")
	sugar.Error("ошибка")          
	logger.Warn("warning...")      
	var apiKey string
	logger.Info("api_key=" + apiKey) 
}

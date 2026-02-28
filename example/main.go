// main.go — расширенный тестовый файл для линтера logcheck (log + zap)
package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	// Инициализация zap (для примера — development режим)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugar := logger.Sugar() // SugaredLogger — самый популярный вариант

	// 1. Uppercase в начале
	log.Printf("User logged in successfully")

	// 2. Русский текст
	sugar.Info("Пользователь вошёл в систему")

	// 3. Специальные символы + эмодзи
	sugar.Warn("Ошибка: токен = abc123!@# 😊")

	// 4. Sensitive данные
	sugar.Error("Failed to authenticate user with password: secret123")

	// 5. Uppercase + русский
	sugar.Infof("Добро пожаловать в приложение!")

	// 6. Специальные символы + эмодзи
	sugar.Warnf("Invalid token detected: %s 😡", "tok_!@#")

	// 7. Sensitive в сообщении
	sugar.Errorw("Login failed",
		"user", "alice",
		"password", "qwerty123", // должно ругаться на sensitive
	)

	// 8. Нормальный английский
	sugar.Info("Request processed", "duration_ms", 42)

	// 9. w-вариант (structured) с потенциально плохими данными
	sugar.Warnw("Config loaded", "api_key", "sk_live_abc123...")

	// 10. Uppercase начало
	logger.Info("User session started successfully")

	// 11. Русский в сообщении
	logger.Warn("Ошибка авторизации пользователя", zap.String("user_id", "123"))

	// 12. Sensitive в поле
	logger.Error("Authentication error",
		zap.String("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
	)

	// 13. Нормальный
	logger.Debug("Health check OK", zap.Int("status", 200))

	fmt.Println("Это не лог → игнорировать")
	sugar.Infow("ok", "field", 42)                           
	logger.Fatal("Critical!", zap.Error(fmt.Errorf("boom"))) 
}

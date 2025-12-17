package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware - логирование запросов
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Обработка запроса
		c.Next()

		// Логирование после обработки
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf("[%s] %s %s %d %v %s",
			method,
			path,
			clientIP,
			statusCode,
			duration,
			c.Errors.String(),
		)

		// Логирование ошибок отдельно
		if len(c.Errors) > 0 {
			log.Printf("ERROR: %s %s - %v", method, path, c.Errors)
		}
	}
}

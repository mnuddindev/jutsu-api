package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	appLogger "github.com/mnuddindev/jutsu-api/internal/infrastructure/logger"
)

// SetupRecover configures recovery middleware
func SetupRecover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
			appLogger.Error("Panic recovered",
				zap.Error(err),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
		},
	})
}

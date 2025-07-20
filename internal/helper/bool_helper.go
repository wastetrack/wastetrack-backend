package helper

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseBoolQuery(ctx *fiber.Ctx, key string) *bool {
	value := ctx.Query(key)
	if value == "" {
		return nil // Parameter not provided
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes":
		result := true
		return &result
	case "false", "0", "no":
		result := false
		return &result
	default:
		return nil // Invalid value, treat as not provided
	}
}

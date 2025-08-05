package helper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ParseBody(ctx *fiber.Ctx, dest interface{}) error {
	err := ctx.BodyParser(dest)
	if err == nil {
		return nil
	}

	var unmarshalTypeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	switch {
	case errors.As(err, &unmarshalTypeErr):
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("Field '%s' expects %s, but got %s",
				unmarshalTypeErr.Field,
				unmarshalTypeErr.Type.String(),
				unmarshalTypeErr.Value),
		)

	case errors.As(err, &syntaxErr):
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("Malformed JSON at offset %d", syntaxErr.Offset))

	default:
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
	}
}

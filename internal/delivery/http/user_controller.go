package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/usecase"
)

type UserController struct {
	Log         *logrus.Logger
	UserUsecase *usecase.UserUseCase
}

func NewUserController(userUsecase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:         logger,
		UserUsecase: userUsecase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	request := new(model.VerifyEmailRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.VerifyEmail(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to verify email: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUsecase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Get(ctx *fiber.Ctx) error {
	request := &model.GetUserRequest{
		ID: ctx.Params("id"),
	}
	response, err := c.UserUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUsecase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

func (c *UserController) LogoutAllDevices(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	err := c.UserUsecase.LogoutAllDevices(ctx.UserContext(), auth.ID)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to logout user from all devices")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *UserController) ResendVerification(ctx *fiber.Ctx) error {
	request := new(model.ResendVerificationRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	err = c.UserUsecase.ResendVerification(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to resend verification: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[string]{
		Data: "success",
	})
}

func (c *UserController) ForgotPassword(ctx *fiber.Ctx) error {
	request := new(model.ForgotPasswordRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	err = c.UserUsecase.ForgotPassword(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to process forgot password: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[string]{
		Data: "success",
	})
}

func (c *UserController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	err = c.UserUsecase.ResetPassword(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to reset password: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[string]{
		Data: "success",
	})
}

func (c *UserController) RefreshToken(ctx *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.RefreshToken(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to refresh token: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

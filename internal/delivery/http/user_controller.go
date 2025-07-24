package http

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/delivery/http/middleware"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
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
	// Get raw body first
	body := ctx.Body()
	c.Log.Infof("Raw request body: %s", string(body))

	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	// Debug parsed values
	c.Log.Infof("Parsed IsAcceptingCustomer: %v", request.IsAcceptingCustomer)

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

func (c *UserController) List(ctx *fiber.Ctx) error {
	// Parse latitude and longitude from query parameters
	var latitude, longitude *float64
	var radiusMeters *int

	if latStr := ctx.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			latitude = &lat
		} else {
			c.Log.Warnf("Invalid latitude parameter: %s", latStr)
			return fiber.NewError(fiber.StatusBadRequest, "Invalid latitude parameter")
		}
	}

	if lngStr := ctx.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			longitude = &lng
		} else {
			c.Log.Warnf("Invalid longitude parameter: %s", lngStr)
			return fiber.NewError(fiber.StatusBadRequest, "Invalid longitude parameter")
		}
	}

	// Parse optional radius parameter (in meters)
	if radiusStr := ctx.Query("radius_meters"); radiusStr != "" {
		if radius, err := strconv.Atoi(radiusStr); err == nil && radius > 0 {
			radiusMeters = &radius
		} else {
			c.Log.Warnf("Invalid radius_meters parameter: %s", radiusStr)
			return fiber.NewError(fiber.StatusBadRequest, "Invalid radius_meters parameter")
		}
	}

	request := &model.SearchUserRequest{
		Username:            ctx.Query("username"),
		Email:               ctx.Query("email"),
		Role:                ctx.Query("role"),
		Institution:         ctx.Query("institution"),
		Address:             ctx.Query("address"),
		City:                ctx.Query("city"),
		Province:            ctx.Query("province"),
		IsAcceptingCustomer: helper.ParseBoolQuery(ctx, "is_accepting_customer"),
		Latitude:            latitude,
		Longitude:           longitude,
		RadiusMeters:        radiusMeters,
		Page:                ctx.QueryInt("page"),
		Size:                ctx.QueryInt("size"),
	}

	responses, total, err := c.UserUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search users with profiles")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.UserListResponse]{
		Data:   responses,
		Paging: paging,
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

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// Parse the request body
	request := new(model.UpdateUserRequest)
	request.ID = ctx.Params("id") // Get user ID from URL params
	request.UserID = auth.ID      // Set the authenticated user's ID for authorization

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	// Call the use case
	response, err := c.UserUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update user: %v", err)
		return err
	}

	// Return the correct response type
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	request := new(model.LogoutUserRequest)
	auth := middleware.GetUser(ctx)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	request.ID = auth.ID

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

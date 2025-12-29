package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type PaginatedResponse struct {
	Status string         `json:"status"`
	Message string        `json:"message"`
	Data   interface{}    `json:"data"`
	Meta   PaginationMeta `json:"meta"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:  "error",
		Message: message,
	})
}

func Paginated(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func SuccessWithMeta(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
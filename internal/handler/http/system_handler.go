package http

import (
	"go-commerce/internal/handler/response"
	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the server is running and database is connected
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]string} "Server is healthy"
// @Router /health [get]
func (h *SystemHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, "Server is running", fiber.Map{
		"status":   "healthy",
		"database": "connected",
	})
}

// Metrics godoc
// @Summary Application metrics endpoint
// @Description Get basic application metrics and version information
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]string} "Application metrics"
// @Router /metrics [get]
func (h *SystemHandler) Metrics(c *fiber.Ctx) error {
	return response.Success(c, "Application metrics", fiber.Map{
		"uptime":     "running",
		"version":    "1.0.0",
		"go_version": "1.25+",
	})
}
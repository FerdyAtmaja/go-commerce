package http

import (
	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type PaymentIntentHandler struct {
	paymentIntentUC domain.PaymentIntentUsecase
}

func NewPaymentIntentHandler(paymentIntentUC domain.PaymentIntentUsecase) *PaymentIntentHandler {
	return &PaymentIntentHandler{
		paymentIntentUC: paymentIntentUC,
	}
}

type CreatePaymentIntentRequest struct {
	Method string `json:"method" validate:"required,oneof=bank_transfer credit_card e_wallet"`
}

// @Summary Create Payment Intent
// @Description Create payment intent for transaction
// @Tags Payment Intent
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body CreatePaymentIntentRequest true "Payment method"
// @Success 201 {object} response.Response{data=domain.PaymentIntent}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /transactions/{id}/pay [post]
func (h *PaymentIntentHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	trxID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	var req CreatePaymentIntentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	intent, err := h.paymentIntentUC.CreatePaymentIntent(uint(trxID), req.Method)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment intent created successfully", intent)
}

// @Summary Simulate Payment Success (Admin Only)
// @Description Simulate payment success callback - Admin only endpoint
// @Tags Payment Simulation (Admin)
// @Accept json
// @Produce json
// @Param intentId path int true "Payment Intent ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /admin/payments/{intentId}/simulate-success [put]
func (h *PaymentIntentHandler) SimulatePaymentSuccess(c *fiber.Ctx) error {
	intentID, err := strconv.ParseUint(c.Params("intentId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid payment intent ID")
	}

	err = h.paymentIntentUC.ProcessPaymentSuccess(uint(intentID))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment processed successfully", nil)
}

// @Summary Simulate Payment Failed (Admin Only)
// @Description Simulate payment failed callback - Admin only endpoint
// @Tags Payment Simulation (Admin)
// @Accept json
// @Produce json
// @Param intentId path int true "Payment Intent ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /admin/payments/{intentId}/simulate-failed [put]
func (h *PaymentIntentHandler) SimulatePaymentFailed(c *fiber.Ctx) error {
	intentID, err := strconv.ParseUint(c.Params("intentId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid payment intent ID")
	}

	err = h.paymentIntentUC.ProcessPaymentFailed(uint(intentID))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment failed processed", nil)
}

type PaymentCallbackRequest struct {
	Status     string `json:"status" validate:"required,oneof=success failed"`
	GatewayRef string `json:"gateway_ref"`
	PaidAt     string `json:"paid_at"`
}

// @Summary Payment Gateway Callback
// @Description Receive payment status from payment gateway - updates payment intent
// @Tags Payment Gateway Callback
// @Accept json
// @Produce json
// @Param intentId path int true "Payment Intent ID"
// @Param request body PaymentCallbackRequest true "Payment status from gateway"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /callbacks/payments/{intentId} [post]
func (h *PaymentIntentHandler) OnPaymentCallback(c *fiber.Ctx) error {
	intentID, err := strconv.ParseUint(c.Params("intentId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid payment intent ID")
	}

	var req PaymentCallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Route to appropriate handler based on status
	switch req.Status {
	case "success":
		err = h.paymentIntentUC.ProcessPaymentSuccess(uint(intentID))
	case "failed":
		err = h.paymentIntentUC.ProcessPaymentFailed(uint(intentID))
	default:
		return response.BadRequest(c, "Invalid payment status")
	}

	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment callback processed", nil)
}
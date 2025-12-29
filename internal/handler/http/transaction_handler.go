package http

import (
	"go-commerce/internal/domain"
	"go-commerce/internal/usecase"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/handler/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionUsecase   *usecase.TransactionUsecase
	paymentIntentUsecase domain.PaymentIntentUsecase
}

func NewTransactionHandler(transactionUsecase *usecase.TransactionUsecase, paymentIntentUsecase domain.PaymentIntentUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase:   transactionUsecase,
		paymentIntentUsecase: paymentIntentUsecase,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction (Authenticated User)
// @Description Create a new transaction with multiple items (atomic operation). Requires authentication.
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateTransactionRequest true "Transaction creation request"
// @Success 200 {object} response.Response{data=domain.Transaction} "Transaction created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req domain.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	transaction, err := h.transactionUsecase.CreateTransaction(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Transaction created successfully", transaction)
}

// GetTransactionByID godoc
// @Summary Get transaction by ID (Authenticated User)
// @Description Get a specific transaction by ID (ownership validation). Only transaction owner can access.
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response{data=domain.Transaction} "Transaction retrieved successfully"
// @Failure 400 {object} response.Response "Invalid transaction ID"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Transaction not found"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	transaction, err := h.transactionUsecase.GetTransactionByID(userID, transactionID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Transaction retrieved successfully", transaction)
}

// GetMyTransactions godoc
// @Summary Get my transactions (Authenticated User)
// @Description Get current user's transaction history with pagination. Requires authentication.
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Transaction} "Transactions retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /transactions/my [get]
func (h *TransactionHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	transactions, total, err := h.transactionUsecase.GetMyTransactions(userID, page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: int((total + int64(limit) - 1) / int64(limit)),
	}

	return response.SuccessWithMeta(c, "Transactions retrieved successfully", transactions, meta)
}



// OnPaymentPaid godoc
// @Summary Payment callback - paid (System)
// @Description Payment gateway callback when payment is successful. System use only.
// @Tags Transactions - Payment Callbacks
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Payment processed successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Router /callbacks/payment/{id}/paid [post]
func (h *TransactionHandler) OnPaymentPaid(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.OnPaymentPaid(transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment processed successfully", nil)
}

// OnPaymentFailed godoc
// @Summary Payment callback - failed (System)
// @Description Payment gateway callback when payment fails. System use only.
// @Tags Transactions - Payment Callbacks
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Payment failure processed successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Router /callbacks/payment/{id}/failed [post]
func (h *TransactionHandler) OnPaymentFailed(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.OnPaymentFailed(transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Payment failure processed successfully", nil)
}

// ProcessOrder godoc
// @Summary Process order (Seller)
// @Description Seller processes order after payment is confirmed. Requires authentication.
// @Tags Transactions - Seller Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Order processed successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /seller/transactions/{id}/process [put]
func (h *TransactionHandler) ProcessOrder(c *fiber.Ctx) error {
	sellerID := middleware.GetUserID(c)
	
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.ProcessOrder(sellerID, transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Order processed successfully", nil)
}

// ShipOrder godoc
// @Summary Ship order (Seller)
// @Description Seller ships processed order. Requires authentication.
// @Tags Transactions - Seller Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Order shipped successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /seller/transactions/{id}/ship [put]
func (h *TransactionHandler) ShipOrder(c *fiber.Ctx) error {
	sellerID := middleware.GetUserID(c)
	
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.ShipOrder(sellerID, transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Order shipped successfully", nil)
}

// ConfirmDelivered godoc
// @Summary Confirm delivery (Buyer)
// @Description Buyer confirms order has been delivered. Requires authentication.
// @Tags Transactions - Buyer Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Delivery confirmed successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /transactions/{id}/confirm-delivery [put]
func (h *TransactionHandler) ConfirmDelivered(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.ConfirmDelivered(userID, transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Delivery confirmed successfully", nil)
}

// CancelTransaction godoc
// @Summary Cancel transaction (Buyer)
// @Description Buyer cancels unpaid transaction. Requires authentication.
// @Tags Transactions - Buyer Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Transaction cancelled successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /transactions/{id}/cancel [put]
func (h *TransactionHandler) CancelTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	// Cancel transaction
	err = h.transactionUsecase.CancelTransaction(userID, transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Expire any active payment intents for this transaction
	// This prevents payment success after cancellation (Koreksi #3)
	h.paymentIntentUsecase.ExpireIntentsByTrxID(uint(transactionID))

	return response.Success(c, "Transaction cancelled successfully", nil)
}

// RefundTransaction godoc
// @Summary Refund transaction (Admin)
// @Description Admin refunds paid transaction. Admin access required.
// @Tags Transactions - Admin Operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Response "Transaction refunded successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /admin/transactions/{id}/refund [put]
func (h *TransactionHandler) RefundTransaction(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	err = h.transactionUsecase.RefundTransaction(transactionID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Transaction refunded successfully", nil)
}
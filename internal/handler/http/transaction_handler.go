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
	transactionUsecase *usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
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
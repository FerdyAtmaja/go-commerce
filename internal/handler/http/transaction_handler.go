package http

import (
	"go-commerce/internal/domain"
	"go-commerce/internal/usecase"
	"go-commerce/internal/handler/response"
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

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint64)

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

func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint64)
	
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

func (h *TransactionHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint64)

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
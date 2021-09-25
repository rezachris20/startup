package handler

import (
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transactions"
	"bwastartup/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionsHandler struct {
	service transactions.Service
	paymentService payment.Service
}

func NewTransactionHandler(service transactions.Service,paymentService payment.Service) *transactionsHandler{
	return &transactionsHandler{service,paymentService}
}

func (h *transactionsHandler) GetCampaignTransactions(c *gin.Context){
	var input transactions.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign transactions",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactionsData,err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	response := helper.APIResponse("Campaign's transactions",http.StatusOK,"success",transactions.FormatCampaignTransactions(transactionsData))
	c.JSON(http.StatusOK,response)
}

func (h *transactionsHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactionsData, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	response := helper.APIResponse("Campaign's transactions",http.StatusOK,"success",transactions.FormatUserTransactions(transactionsData))
	c.JSON(http.StatusOK,response)
}

func (h *transactionsHandler) CreateTransaction(c *gin.Context) {
	var input transactions.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Failed to created transaction",http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity,response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newTransaction,err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to created transaction",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	response := helper.APIResponse("Success to created transaction",http.StatusOK,"success",transactions.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK,response)
}

func (h *transactionsHandler) GetNotification(c *gin.Context) {
	var input transactions.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		res := helper.APIResponse("Failed to process notification",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,res)
		return
	}

	err = h.paymentService.ProcessPayment(input)
	if err != nil {
		res := helper.APIResponse("Failed to process notification",http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest,res)
		return
	}

	c.JSON(http.StatusOK,input)
}
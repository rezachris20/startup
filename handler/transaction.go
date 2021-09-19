package handler

import (
	"bwastartup/helper"
	"bwastartup/transactions"
	"bwastartup/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionsHandler struct {
	service transactions.Service
}

func NewTransactionHandler(service transactions.Service) *transactionsHandler{
	return &transactionsHandler{service}
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
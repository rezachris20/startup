package payment

import (
	"bwastartup/campaign"
	"bwastartup/transactions"
	"bwastartup/user"
	"github.com/veritrans/go-midtrans"
	"strconv"
)

type service struct {
	transactionRepository transactions.Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transactions.TransactionNotificationInput) error
}

func NewService(transactionRepository transactions.Repository, campaignRepository campaign.Repository) *service {
	return &service{transactionRepository, campaignRepository}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-vDOcZSJt6OnG8qrwC0nXfRbp"
	midclient.ClientKey = "SB-Mid-client-wAcfHRiKrN7bmupJ"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},

		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

func (s *service) ProcessPayment(input transactions.TransactionNotificationInput) error {
	transactionID, _ := strconv.Atoi(input.OrderID)

	transactions,err := s.transactionRepository.GetByID(transactionID)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transactions.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transactions.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transactions.Status = "cancel"
	}

	updatedTrasaction,err := s.transactionRepository.Update(transactions)
	if err != nil {
		return err
	}

	campaign,err := s.campaignRepository.FindByID(updatedTrasaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTrasaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTrasaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}

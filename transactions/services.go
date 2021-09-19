package transactions

import (
	"bwastartup/campaign"
	"errors"
)

type service struct {
	repository Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput)([]Transaction,error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service)GetTransactionsByCampaignID(input GetCampaignTransactionsInput)([]Transaction,error){
	//get campaign
	campaign,err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{},err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{},errors.New("Not an owner of the campaign")
	}
	//check campaign.userid != user_id_yang_melakukan_request
	transactions,err := s.repository.GetCampaignID(input.ID)
	if err != nil {
		return transactions,err
	}
	return  transactions,nil
}

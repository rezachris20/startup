package transactions

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetCampaignID(campaignID int)([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository)GetCampaignID(campaignID int)([]Transaction, error){
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?",campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions,err
	}

	return transactions,nil
}
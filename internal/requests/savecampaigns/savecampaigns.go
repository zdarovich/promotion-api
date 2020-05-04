package savecampaigns

import (
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignhelper"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	"strconv"
	"time"
)

type (
	// SaveCampaigns struct
	SaveCampaigns struct {
		CampaignRepository campaign.IRepository
		CampaignHelper     campaignhelper.ICampaignHelper
		Configuration      *config.Configuration
		InputParameters    inputParameters
	}
	// requestParams the parameters that can be used for searching
	inputParameters struct {
		ID                      int       `json:"id"`
		StartDate               time.Time `json:"start_date"`
		EndDate                 time.Time `json:"end_date"`
		Name                    string    `json:"name"`
		WarehouseID             int       `json:"warehouse_id"`
		PurchasedAmount         int       `json:"purchased_amount"`
		PurchasedProdgroupID    int       `json:"purchased_prodgroup_id"`
		PurchaseTotalValue      float64   `json:"purchase_total_value"`
		AwardLowestPricedItem   int       `json:"award_lowest_priced_item"`
		SpecialPrice            float64   `json:"special_price"`
		PercentageOff           int       `json:"percentage_off"`
		SumOff                  float64   `json:"sum_off"`
		AwardedProdgroupID      int       `json:"awarded_prodgroup_id"`
		PercentageOffAllItems   int       `json:"percentage_off_all_items"`
		SumOffEntirePurchase    float64   `json:"sum_off_entire_purchase"`
		Rewardpoints            int       `json:"rewardpoints"`
		PercentageOffAnyOneLine int       `json:"percentage_off_any_one_line"`
		Type                    string    `json:"type"`
		Added                   int       `json:"added"`
		Addedby                 string    `json:"addedby"`
		Changed                 int       `json:"changed"`
		Changedby               string    `json:"changedby"`
	}
)

// @Summary Save campaign
// @Description  Save campaign
// @Tags campaign
// @Accept  json
// @Produce  json
// @Param sessionKey formData string true "session key"
// @Param clientCode formData string true "client code"
// @Param request formData string true "client code"
// @Param campaignID formData string true "campaign IDs - (1,4,7)"
// @Param startDate formData string true "2006-01-02"
// @Param endDate formData string true "2006-01-02"
// @Param name formData string true "test"
// @Param warehouseID formData string true "1"
// @Param purchasedAmount formData string true "1"
// @Param purchaseTotalValue formData string true "1"
// @Param purchasedProductGroupID formData string true "1"
// @Param specialPrice formData string true "1"
// @Param percentageOff formData string true "1"
// @Param sumOff formData string true "1"
// @Param awardedProductGroupID formData string true "1"
// @Param percentageOffMatchingItems formData string true "1"
// @Param type formData string true "1"
// @Success 200 "Created"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router / [POST]
func (saveCampaigns *SaveCampaigns) Handle(context root.IGinContext) (*response.Data, error) {

	err := saveCampaigns.validate(context)

	if err != nil {
		return nil, err
	}

	var totalRecordsCount int = 0
	var recordsCount int = 0
	var records interface{}

	c := campaign.Campaign{
		saveCampaigns.InputParameters.ID,
		saveCampaigns.InputParameters.StartDate,
		saveCampaigns.InputParameters.EndDate,
		saveCampaigns.InputParameters.Name,
		saveCampaigns.InputParameters.WarehouseID,
		saveCampaigns.InputParameters.PurchasedAmount,
		saveCampaigns.InputParameters.PurchasedProdgroupID,
		saveCampaigns.InputParameters.PurchaseTotalValue,
		saveCampaigns.InputParameters.AwardLowestPricedItem,
		saveCampaigns.InputParameters.SpecialPrice,
		saveCampaigns.InputParameters.PercentageOff,
		saveCampaigns.InputParameters.SumOff,
		saveCampaigns.InputParameters.AwardedProdgroupID,
		saveCampaigns.InputParameters.PercentageOffAllItems,
		saveCampaigns.InputParameters.SumOffEntirePurchase,
		saveCampaigns.InputParameters.Rewardpoints,
		saveCampaigns.InputParameters.PercentageOffAnyOneLine,
		saveCampaigns.InputParameters.Type,
		saveCampaigns.InputParameters.Added,
		saveCampaigns.InputParameters.Addedby,
		saveCampaigns.InputParameters.Changed,
		saveCampaigns.InputParameters.Changedby,
	}
	err = saveCampaigns.CampaignHelper.Validate(&c)
	if err != nil {
		return nil, err
	}
	err = saveCampaigns.CampaignRepository.SaveCampaigns(c)

	if err != nil {
		return nil, err
	}

	totalRecordsCount, err = saveCampaigns.CampaignRepository.GetCampaignsCount(
		0,
		"",
	)
	if err != nil {
		return nil, err
	}
	recordsCount = 0

	return &response.Data{
		Total:           totalRecordsCount,
		TotalInResponse: recordsCount,
		Records:         records,
	}, nil
}

// New return configured struct
func New(configuration *config.Configuration) root.IRoot {

	return &SaveCampaigns{
		CampaignRepository: campaign.New(configuration),
		CampaignHelper:     campaignhelper.New(configuration),
		Configuration:      configuration,
	}
}

// validate checks if the required parameters have been set
func (saveCampaigns *SaveCampaigns) validate(context root.IGinContext) error {

	inputParameters := inputParameters{}
	inputParameters.ID, _ = strconv.Atoi(context.PostForm("campaignID"))
	inputParameters.StartDate, _ = time.Parse("2006-01-02", context.PostForm("startDate"))
	inputParameters.EndDate, _ = time.Parse("2006-01-02", context.PostForm("endDate"))
	inputParameters.Name = context.PostForm("name")
	inputParameters.WarehouseID, _ = strconv.Atoi(context.PostForm("warehouseID"))
	inputParameters.PurchasedAmount, _ = strconv.Atoi(context.PostForm("purchasedAmount"))
	inputParameters.PurchaseTotalValue, _ = strconv.ParseFloat(context.PostForm("purchaseTotalValue"), 64)
	inputParameters.AwardLowestPricedItem, _ = strconv.Atoi(context.PostForm("purchasedProductGroupID"))
	inputParameters.SpecialPrice, _ = strconv.ParseFloat(context.PostForm("specialPrice"), 64)
	inputParameters.PercentageOff, _ = strconv.Atoi(context.PostForm("percentageOff"))
	inputParameters.SumOff, _ = strconv.ParseFloat(context.PostForm("sumOff"), 64)
	inputParameters.AwardedProdgroupID, _ = strconv.Atoi(context.PostForm("awardedProductGroupID"))
	inputParameters.PercentageOffAllItems, _ = strconv.Atoi(context.PostForm("percentageOffMatchingItems"))
	inputParameters.Type = context.PostForm("type")

	saveCampaigns.InputParameters = inputParameters
	return nil
}

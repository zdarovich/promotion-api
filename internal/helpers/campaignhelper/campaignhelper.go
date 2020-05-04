package campaignhelper

import (
	"errors"
	"fmt"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignvalidator"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
)

type (
	// CampaignHelper struct
	CampaignHelper struct {
		Configuration      *config.Configuration
		CampaignRepository campaign.IRepository
	}
	// ICampaignHelper interface
	ICampaignHelper interface {
		MapToArray(campaigns []campaign.Campaign) interface{}
		Validate(campaigns *campaign.Campaign) error
	}
	// record structure of the output record
	Record struct {
		CampaignID                                       int           `json:"campaignID"`
		StartDate                                        string        `json:"startDate"`
		EndDate                                          string        `json:"endDate"`
		Name                                             string        `json:"name"`
		Type                                             string        `json:"type"`
		WarehouseID                                      int           `json:"warehouseID"`
		PurchasedProductGroupID                          int           `json:"purchasedProductGroupID"`
		PurchasedBrandID                                 int           `json:"purchasedBrandID"`
		PurchasedAmount                                  int           `json:"purchasedAmount"`
		PurchaseTotalValue                               string        `json:"purchaseTotalValue"`
		RewardPoints                                     int           `json:"rewardPoints"`
		PercentageOffEntirePurchase                      int           `json:"percentageOffEntirePurchase"`
		SumOffEntirePurchase                             string        `json:"sumOffEntirePurchase"`
		SpecialPrice                                     string        `json:"specialPrice"`
		AwardedProductGroupID                            int           `json:"awardedProductGroupID"`
		AwardedBrandID                                   int           `json:"awardedBrandID"`
		LowestPriceItemIsAwarded                         int           `json:"lowestPriceItemIsAwarded"`
		PercentageOFF                                    int           `json:"percentageOFF"`
		SumOFF                                           string        `json:"sumOFF"`
		DiscountForOneLine                               int           `json:"discountForOneLine"`
		Added                                            int           `json:"added"`
		LastModified                                     int           `json:"lastModified"`
		RequiredCouponID                                 string        `json:"requiredCouponID"`
		RequiredCouponCode                               string        `json:"requiredCouponCode"`
		PurchasedProducts                                []interface{} `json:"purchasedProducts"`
		AwardedProducts                                  []interface{} `json:"awardedProducts"`
		ExcludedProducts                                 []interface{} `json:"excludedProducts"`
		PercentageOffExcludedProducts                    []interface{} `json:"percentageOffExcludedProducts"`
		SumOffExcludedProducts                           []interface{} `json:"sumOffExcludedProducts"`
		SumOffIncludedProducts                           []interface{} `json:"sumOffIncludedProducts"`
		AwardedAmount                                    int           `json:"awardedAmount"`
		PurchasedProductCategoryID                       int           `json:"purchasedProductCategoryID"`
		AwardedProductCategoryID                         int           `json:"awardedProductCategoryID"`
		MaximumPointsDiscount                            int           `json:"maximumPointsDiscount"`
		CustomerCanUseOnlyOnce                           int           `json:"customerCanUseOnlyOnce"`
		PriceAtLeast                                     int           `json:"priceAtLeast"`
		PriceAtMost                                      int           `json:"priceAtMost"`
		RequiresManagerOverride                          int           `json:"requiresManagerOverride"`
		SumOffMatchingItems                              int           `json:"sumOffMatchingItems"`
		PercentageOffMatchingItems                       int           `json:"percentageOffMatchingItems"`
		ExcludeDiscountedFromPercentageOffEntirePurchase int           `json:"excludeDiscountedFromPercentageOffEntirePurchase"`
		ReasonID                                         int           `json:"reasonID"`
		SpecialUnitPrice                                 int           `json:"specialUnitPrice"`
		MaxItemsWithSpecialUnitPrice                     int           `json:"maxItemsWithSpecialUnitPrice"`
		RedemptionLimit                                  int           `json:"redemptionLimit"`
		StoreGroup                                       string        `json:"storeGroup"`
		CanBeAppliedManuallyMultipleTimes                int           `json:"canBeAppliedManuallyMultipleTimes"`
	}
)

// New returns configured product file helper
func New(configuration *config.Configuration) ICampaignHelper {

	return &CampaignHelper{
		Configuration:      configuration,
		CampaignRepository: campaign.New(configuration),
	}
}

// MapToArray maps database files to records
func (p *CampaignHelper) MapToArray(campaigns []campaign.Campaign) interface{} {

	records := make([]Record, 0)

	for _, c := range campaigns {
		records = append(records, Record{
			CampaignID:                  c.ID,
			StartDate:                   c.StartDate.String(),
			EndDate:                     c.EndDate.String(),
			Name:                        c.Name,
			WarehouseID:                 c.WarehouseID,
			PurchasedAmount:             c.PurchasedAmount,
			PurchasedProductGroupID:     c.PurchasedProdgroupID,
			PurchaseTotalValue:          fmt.Sprintf("%f", c.PurchaseTotalValue),
			AwardedAmount:               c.AwardLowestPricedItem,
			SpecialPrice:                fmt.Sprintf("%f", c.SpecialPrice),
			PercentageOFF:               c.PercentageOff,
			SumOFF:                      fmt.Sprintf("%f", c.SumOff),
			AwardedProductGroupID:       c.AwardedProdgroupID,
			PercentageOffMatchingItems:  c.PercentageOffAllItems,
			SumOffEntirePurchase:        fmt.Sprintf("%f", c.SumOffEntirePurchase),
			RewardPoints:                c.Rewardpoints,
			PercentageOffEntirePurchase: c.PercentageOffAnyOneLine,
			Type:                        c.Type,
			Added:                       c.Added,
		})
	}

	return records
}

func (p *CampaignHelper) Validate(campaign *campaign.Campaign) error {
	if campaignvalidator.IsOnlyOneAwardIsDefined(campaign) {
		return errors.New("1048")
	} else if campaignvalidator.IsLowestPriceItem(campaign) {
		return errors.New("1119")

	} else if campaignvalidator.IsOnlyOneRequirement(campaign) {
		return errors.New("1141")

	} else if campaignvalidator.IsMultiBuyRequirement(campaign) {
		return errors.New("1123")

	} else if campaignvalidator.IsSpecialPriceRequirement(campaign) {
		return errors.New("1122")
	}
	return nil
}

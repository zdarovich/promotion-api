package campaignhelper

import (
	"errors"
	"fmt"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
	"github.com/zdarovich/promotion-api/internal/repositories/attributes"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	configurationRepo "github.com/zdarovich/promotion-api/internal/repositories/config"
	"reflect"
	"strings"
	"time"
)

type (
	// CampaignHelper struct
	CampaignHelper struct {
		Configuration      *config.Configuration
		CampaignRepository campaign.IRepository
		ConfigRepository   configurationRepo.IRepository
	}
	// ICampaignHelper interface
	ICampaignHelper interface {
		MapToArray(cs []campaign.Campaign, attrs map[int][]*attributes.Attribute) ([]RecordOutput, error)
		MapToOutput(records []Record) ([]RecordOutput, error)
		Validate(attrs *Record) error
	}
	// record structure of the output record
	RecordOutput struct {
		CampaignID                                           int       `json:"campaignID"`
		StartDate                                            time.Time `json:"startDate"`
		EndDate                                              time.Time `json:"endDate"`
		Name                                                 string    `json:"name"`
		Type                                                 string    `json:"type"`
		WarehouseID                                          int       `json:"warehouseID"`
		AwardedProductGroupID                                int       `json:"awardedProductGroupID"`
		AwardedBrandID                                       int       `json:"awardedBrandID"`
		LowestPriceItemIsAwarded                             int       `json:"lowestPriceItemIsAwarded"`
		PercentageOFF                                        float64   `json:"percentageOFF"`
		SumOFF                                               float64   `json:"sumOFF"`
		DiscountForOneLine                                   int       `json:"discountForOneLine"`
		RequiredCouponID                                     string    `json:"requiredCouponID"`
		RequiredCouponCode                                   string    `json:"requiredCouponCode"`
		PurchasedProducts                                    string    `json:"purchasedProducts"`
		AwardedProducts                                      string    `json:"awardedProducts"`
		ExcludedProducts                                     string    `json:"excludedProducts"`
		PercentageOffExcludedProducts                        string    `json:"percentageOffExcludedProducts"`
		PercentageOffIncludedProducts                        string    `json:"percentageOffIncludedProducts"`
		PurchasedProductSubsidies                            string    `json:"purchasedProductSubsidies"`
		SumOffExcludedProducts                               string    `json:"sumOffExcludedProducts"`
		SumOffIncludedProducts                               string    `json:"sumOffIncludedProducts"`
		AwardedProductSubsidies                              string    `json:"awardedProductSubsidies"`
		StoreRegionIDs                                       string    `json:"storeRegionIDs"`
		CustomerGroupIDs                                     string    `json:"customerGroupIDs"`
		AwardedAmount                                        int       `json:"awardedAmount"`
		PurchasedProductCategoryID                           int       `json:"purchasedProductCategoryID"`
		AwardedProductCategoryID                             int       `json:"awardedProductCategoryID"`
		MaximumPointsDiscount                                int       `json:"maximumPointsDiscount"`
		CustomerCanUseOnlyOnce                               int       `json:"customerCanUseOnlyOnce"`
		PriceAtLeast                                         int       `json:"priceAtLeast"`
		PriceAtMost                                          int       `json:"priceAtMost"`
		RequiresManagerOverride                              int       `json:"requiresManagerOverride"`
		SumOffMatchingItems                                  int       `json:"sumOffMatchingItems"`
		PercentageOffMatchingItems                           int       `json:"percentageOffMatchingItems"`
		ExcludeDiscountedFromPercentageOffEntirePurchase     int       `json:"excludeDiscountedFromPercentageOffEntirePurchase"`
		ExcludePromotionItemsFromPercentageOffEntirePurchase int       `json:"excludePromotionItemsFromPercentageOffEntirePurchase"`
		ReasonID                                             int       `json:"reasonID"`
		SpecialUnitPrice                                     int       `json:"specialUnitPrice"`
		MaxItemsWithSpecialUnitPrice                         int       `json:"maxItemsWithSpecialUnitPrice"`
		RedemptionLimit                                      int       `json:"redemptionLimit"`
		StoreGroup                                           string    `json:"storeGroup"`
		CanBeAppliedManuallyMultipleTimes                    int       `json:"canBeAppliedManuallyMultipleTimes"`
		PurchasedProductGroupID                              int       `json:"purchasedProductGroupID"`
		PurchasedBrandID                                     int       `json:"purchasedBrandID"`
		PurchasedAmount                                      int       `json:"purchasedAmount"`
		PurchaseTotalValue                                   float64   `json:"purchaseTotalValue"`
		RewardPoints                                         int       `json:"rewardPoints"`
		PercentageOffEntirePurchase                          int       `json:"percentageOffEntirePurchase"`
		SumOffEntirePurchase                                 float64   `json:"sumOffEntirePurchase"`
		SpecialPrice                                         float64   `json:"specialPrice"`
		Added                                                int64     `json:"added"`
		Addedby                                              string    `json:"addedby"`
		Changed                                              int64     `json:"changed"`
		Changedby                                            string    `json:"changedby"`
	}
	Record struct {
		CampaignID                                           int       `json:"campaignID"`
		StartDate                                            time.Time `json:"startDate"`
		EndDate                                              time.Time `json:"endDate"`
		Name                                                 string    `json:"name"`
		Type                                                 string    `json:"type"`
		WarehouseID                                          int       `json:"warehouseID"`
		AwardedProductGroupID                                int       `json:"awardedProductGroupID"`
		AwardedBrandID                                       int       `json:"awardedBrandID"`
		LowestPriceItemIsAwarded                             bool      `json:"lowestPriceItemIsAwarded"`
		PercentageOFF                                        float64   `json:"percentageOFF"`
		SumOFF                                               float64   `json:"sumOFF"`
		DiscountForOneLine                                   int       `json:"discountForOneLine"`
		RequiredCouponID                                     string    `json:"requiredCouponID"`
		RequiredCouponCode                                   string    `json:"requiredCouponCode"`
		PurchasedProducts                                    []string  `json:"purchasedProducts"`
		AwardedProducts                                      []string  `json:"awardedProducts"`
		ExcludedProducts                                     []string  `json:"excludedProducts"`
		PercentageOffExcludedProducts                        []string  `json:"percentageOffExcludedProducts"`
		PercentageOffIncludedProducts                        []string  `json:"percentageOffIncludedProducts"`
		PurchasedProductSubsidies                            []string  `json:"purchasedProductSubsidies"`
		SumOffExcludedProducts                               []string  `json:"sumOffExcludedProducts"`
		SumOffIncludedProducts                               []string  `json:"sumOffIncludedProducts"`
		AwardedProductSubsidies                              []string  `json:"awardedProductSubsidies"`
		StoreRegionIDs                                       []int     `json:"storeRegionIDs"`
		CustomerGroupIDs                                     []int     `json:"customerGroupIDs"`
		AwardedAmount                                        int       `json:"awardedAmount"`
		PurchasedProductCategoryID                           int       `json:"purchasedProductCategoryID"`
		AwardedProductCategoryID                             int       `json:"awardedProductCategoryID"`
		MaximumPointsDiscount                                int       `json:"maximumPointsDiscount"`
		CustomerCanUseOnlyOnce                               bool      `json:"customerCanUseOnlyOnce"`
		PriceAtLeast                                         int       `json:"priceAtLeast"`
		PriceAtMost                                          int       `json:"priceAtMost"`
		RequiresManagerOverride                              bool      `json:"requiresManagerOverride"`
		SumOffMatchingItems                                  int       `json:"sumOffMatchingItems"`
		PercentageOffMatchingItems                           int       `json:"percentageOffMatchingItems"`
		ExcludeDiscountedFromPercentageOffEntirePurchase     bool      `json:"excludeDiscountedFromPercentageOffEntirePurchase"`
		ExcludePromotionItemsFromPercentageOffEntirePurchase bool      `json:"excludePromotionItemsFromPercentageOffEntirePurchase"`
		ReasonID                                             int       `json:"reasonID"`
		SpecialUnitPrice                                     int       `json:"specialUnitPrice"`
		MaxItemsWithSpecialUnitPrice                         int       `json:"maxItemsWithSpecialUnitPrice"`
		RedemptionLimit                                      int       `json:"redemptionLimit"`
		StoreGroup                                           string    `json:"storeGroup"`
		CanBeAppliedManuallyMultipleTimes                    int       `json:"canBeAppliedManuallyMultipleTimes"`
		PurchasedProductGroupID                              int       `json:"purchasedProductGroupID"`
		PurchasedBrandID                                     int       `json:"purchasedBrandID"`
		PurchasedAmount                                      int       `json:"purchasedAmount"`
		PurchaseTotalValue                                   float64   `json:"purchaseTotalValue"`
		RewardPoints                                         int       `json:"rewardPoints"`
		PercentageOffEntirePurchase                          int       `json:"percentageOffEntirePurchase"`
		SumOffEntirePurchase                                 float64   `json:"sumOffEntirePurchase"`
		SpecialPrice                                         float64   `json:"specialPrice"`
		Added                                                int64     `json:"added"`
		Addedby                                              string    `json:"addedby"`
		Changed                                              int64     `json:"changed"`
		Changedby                                            string    `json:"changedby"`
	}
)

// New returns configured product file helper
func New(configuration *config.Configuration) ICampaignHelper {

	return &CampaignHelper{
		Configuration:    configuration,
		ConfigRepository: configurationRepo.New(configuration),
	}
}

func (p *CampaignHelper) MapToArray(cs []campaign.Campaign, attrs map[int][]*attributes.Attribute) ([]RecordOutput, error) {
	result := make([]RecordOutput, len(cs))
	for idx, c := range cs {
		ro := RecordOutput{}
		ro.CampaignID = c.ID
		ro.StartDate = c.StartDate
		ro.EndDate = c.EndDate
		ro.Name = c.Name
		ro.Type = c.Type
		ro.WarehouseID = c.WarehouseID
		ro.AwardedProductGroupID = c.AwardedProdgroupID
		var awardLowestPricedItem int
		if c.AwardLowestPricedItem {
			awardLowestPricedItem = 1
		}
		ro.LowestPriceItemIsAwarded = awardLowestPricedItem
		ro.PercentageOFF = c.PercentageOff
		ro.SumOFF = c.SumOff
		ro.PercentageOffMatchingItems = c.PercentageOffAnyOneLine
		ro.PurchasedProductGroupID = c.PurchasedProdgroupID
		ro.PurchasedAmount = c.PurchasedAmount
		ro.PurchaseTotalValue = c.PurchaseTotalValue
		ro.RewardPoints = c.Rewardpoints
		ro.PercentageOffEntirePurchase = c.PercentageOffAllItems
		ro.SumOffEntirePurchase = c.SumOffEntirePurchase
		ro.SpecialPrice = c.SpecialPrice
		ro.Added = c.Added
		ro.Addedby = c.Addedby
		ro.Changed = c.Changed
		ro.Changedby = c.Changedby

		vals := make(map[string]interface{})
		for _, attr := range attrs[c.ID] {
			switch attr.Type {
			case "text":
				vals[attr.Name] = attr.ValueText
			case "int":
				vals[attr.Name] = attr.ValueInt
			case "double":
				vals[attr.Name] = attr.ValueDouble
			}
		}

		v := reflect.ValueOf(&ro)
		v = reflect.Indirect(v)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			val, ok := vals[field]
			if !ok {
				continue
			}

			value := v.Field(i)
			switch e := val.(type) {
			case int:
				value.SetInt(int64(e))
			case string:
				value.SetString(e)
			case float64:
				value.SetFloat(e)
			}
		}
		result[idx] = ro
	}
	return result, nil
}

// MapToOutput maps database files to records
func (p *CampaignHelper) MapToOutput(records []Record) ([]RecordOutput, error) {

	output := make([]RecordOutput, 0)

	for _, r := range records {
		v := reflect.ValueOf(r)
		v = reflect.Indirect(v)
		t := v.Type()
		vals := make(map[string]reflect.Value)

		for i := 0; i < v.NumField(); i++ {
			val := v.Field(i)
			if val.IsZero() {
				continue
			}
			field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			vals[field] = val
		}
		ro := RecordOutput{}
		v = reflect.ValueOf(&ro)
		v = reflect.Indirect(v)
		t = v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			val, ok := vals[field]
			if !ok {
				continue
			}
			switch val.Kind() {
			case reflect.Slice:
				switch val.Type().String() {
				case "[]int":
					log.Info(val.Interface().([]int))
					v.Field(i).SetString(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(val.Interface().([]int))), ","), "[]"))
				case "[]string":
					v.Field(i).SetString(strings.Join(val.Interface().([]string), ","))
				}

			case reflect.Bool:
				var bitSet int64 = 0
				if val.Interface().(bool) {
					bitSet = 1
				}
				v.Field(i).SetInt(bitSet)
			default:
				v.Field(i).Set(val)
			}
		}
		output = append(output, ro)
	}
	return output, nil
}

func (p *CampaignHelper) Validate(r *Record) error {
	if r == nil {
		return errors.New("record is null")
	}
	c, err := p.ConfigRepository.GetConfigByName("vertical")
	if err != nil {
		log.Error(err)
		return errors.New("1006")
	}

	if !IsTypeValid(r) {
		return errorcodes.New("type", 1014)
	} else if !IsStoreRegionIDsEnabled(r, c) {
		return errors.New("1028")
	} else if !IsCustomerGroupIDsEnabled(r, c) {
		return errors.New("1028")
	} else if !IsStartDateValid(r) {
		return errorcodes.New("startDate", 1014)
	} else if !IsEndDateValid(r) {
		return errorcodes.New("endDate", 1014)
	} else if IsRequiresManagerOverrideAndNotAutomaticOrNotCoupon(r) {
		return errors.New("1076")
	}

	// business requirements
	if IsMultipleSetting(r) {
		return errors.New("1110")
	} else if !IsPurchasedProductGroupIDOrPurchasedProductCategoryIDOrPurchasedProductsAndPurchasedAmount(r) {
		return errors.New("1111")
	} else if IsMultipleProductOptions(r) {
		return errors.New("1112")
	} else if !IsAwardedProductOptionsAndSumOffOrPercentageOff(r) {
		return errors.New("1113")
	} else if IsMultipleAwardOptions(r) {
		return errors.New("1114")
	} else if !IsPercentageExclInclProductsAndPercentageOffEntirePurchase(r) {
		return errors.New("1115")
	} else if !IsSumExclInclProductsAndSumOffEntirePurchase(r) {
		return errors.New("1116")
	} else if !IsPriceAtLeastOrPriceAtMostAndPurchasedAmount(r) {
		return errors.New("1117")
	} else if !IsMaximumPointsDiscountAndRewardPointsAndSumOffEntirePurchase(r) {
		return errors.New("1118")
	} else if !IsLowestPriceItemIsAwardedAndSumOffOrPercentageOff(r) {
		return errors.New("1119")
	} else if !IsSpecialPriceAndPurchasedAmount(r) {
		return errors.New("1122")
	} else if !IsPercentageOffMatchingItemsOrSumOffEntirePurchaseAndPurchasedAmount(r) {
		return errors.New("1123")
	} else if !IsExcludeDiscountedFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(r) {
		return errors.New("1129")
	} else if !IsExcludePromotionItemsFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(r) {
		return errors.New("1182")
	} else if !IsReasonIdEqualsPromotion(r) {
		return errors.New("1131")
	} else if !IsPurchasedProductSubsidiesAndPurchasedProductsAndPercentageOffMatchingItemsOrSumOffMatchingItems(r) {
		return errors.New("1132")
	} else if !IsPurchasedProductSubsidiesLenEqualsPurchasedProductsLen(r) {
		return errors.New("1133")
	} else if !IsAwardedProductSubsidiesLenEqualsAwardedProductsLen(r) {
		return errors.New("1134")
	} else if !IsSpecialUnitPriceAndPurchasedAmount(r) {
		return errors.New("1139")
	} else if !IsMaxItemsWithSpecialUnitPriceEqualsOrBiggerPurchasedAmount(r) {
		return errors.New("1140")
	} else if !IsPurchasedAmountAndPurchasedProductGroupIDOrPurchasedProductCategoryIDOrPurchasedProducts(r) {
		return errors.New("1141")
	} else if !IsRedemptionLimitAndNotPercentageOfEntirePurchaseAndNotRewardPoints(r) {
		return errors.New("1144")
	} else if !IsRedemptionLimitAndMaxItemsWithSpecialUnitPrice(r) {
		return errors.New("1145")
	}
	return nil
}

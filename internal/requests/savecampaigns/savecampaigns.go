package savecampaigns

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignhelper"
	"github.com/zdarovich/promotion-api/internal/log"
	"github.com/zdarovich/promotion-api/internal/repositories/attributes"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	"github.com/zdarovich/promotion-api/internal/repositories/user"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	// SaveCampaigns struct
	SaveCampaigns struct {
		CampaignRepository campaign.IRepository
		AttrsRepository    attributes.IRepository
		CampaignHelper     campaignhelper.ICampaignHelper
		UserRepository     user.IRepository
		Configuration      *config.Configuration
	}

	inputParametersAttrs struct {
		AwardedBrandID                                       int    `json:"awardedBrandID"`
		DiscountForOneLine                                   int    `json:"discountForOneLine"`
		RequiredCouponID                                     string `json:"requiredCouponID"`
		RequiredCouponCode                                   string `json:"requiredCouponCode"`
		PurchasedProducts                                    string `json:"purchasedProducts"`
		AwardedProducts                                      string `json:"awardedProducts"`
		ExcludedProducts                                     string `json:"excludedProducts"`
		PercentageOffExcludedProducts                        string `json:"percentageOffExcludedProducts"`
		PercentageOffIncludedProducts                        string `json:"percentageOffIncludedProducts"`
		SumOffExcludedProducts                               string `json:"sumOffExcludedProducts"`
		SumOffIncludedProducts                               string `json:"sumOffIncludedProducts"`
		AwardedAmount                                        int    `json:"awardedAmount"`
		PurchasedProductCategoryID                           int    `json:"purchasedProductCategoryID"`
		AwardedProductCategoryID                             int    `json:"awardedProductCategoryID"`
		MaximumPointsDiscount                                int    `json:"maximumPointsDiscount"`
		CustomerCanUseOnlyOnce                               int    `json:"customerCanUseOnlyOnce"`
		PriceAtLeast                                         int    `json:"priceAtLeast"`
		PriceAtMost                                          int    `json:"priceAtMost"`
		RequiresManagerOverride                              int    `json:"requiresManagerOverride"`
		SumOffMatchingItems                                  int    `json:"sumOffMatchingItems"`
		ExcludeDiscountedFromPercentageOffEntirePurchase     int    `json:"excludeDiscountedFromPercentageOffEntirePurchase"`
		ExcludePromotionItemsFromPercentageOffEntirePurchase int    `json:"excludePromotionItemsFromPercentageOffEntirePurchase"`
		ReasonID                                             int    `json:"reasonID"`
		SpecialUnitPrice                                     int    `json:"specialUnitPrice"`
		MaxItemsWithSpecialUnitPrice                         int    `json:"maxItemsWithSpecialUnitPrice"`
		RedemptionLimit                                      int    `json:"redemptionLimit"`
		StoreGroup                                           string `json:"storeGroup"`
		CanBeAppliedManuallyMultipleTimes                    int    `json:"canBeAppliedManuallyMultipleTimes"`
		PurchasedBrandID                                     int    `json:"purchasedBrandID"`
		AwardedProductGroupID                                int    `json:"awardedProductGroupID"`
		LowestPriceItemIsAwarded                             int    `json:"lowestPriceItemIsAwarded"`
		PurchasedProductSubsidies                            string `json:"purchasedProductSubsidies"`
		AwardedProductSubsidies                              string `json:"awardedProductSubsidies"`
		StoreRegionIDs                                       string `json:"storeRegionIDs"`
		CustomerGroupIDs                                     string `json:"customerGroupIDs"`
		PercentageOffMatchingItems                           int    `json:"percentageOffMatchingItems"`
		PurchasedProductGroupID                              int    `json:"purchasedProductGroupID"`
		RewardPoints                                         int    `json:"rewardPoints"`
		PercentageOffEntirePurchase                          int    `json:"percentageOffEntirePurchase"`
	}
)

// @Summary Save campaign
// @Description  Save campaign
// @Tags campaign
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param sessionKey formData string true "ERPLY session key"
// @Param clientCode formData string true "ERPLY client code"
// @Param request formData string true "saveCampaign"
// @Param campaignID formData string false "1"
// @Description  startDate - Promotion start date.
// @Param startDate formData string false "2006-01-02"
// @Description  endDate - Promotion end date.
// @Param endDate formData string false "2006-01-02"
// @Description  name - Promotion name. Use either general parameter "name" or one or more of the following parameters if you need to set the names in specific languages. To have multilingual names enabled, please contact our customer support. An error will be returned if you attempt to set a name in a specific language and the multilingual names are not enabled on your account.
// @Param name formData string false "test"
// @Description  warehouseID - Set this field if you want the promotion to be available only in a specific store.Fields "warehouseID", "storeGroup" and "storeRegionIDs" are mutually exclusive: only one restriction can be set at a time, otherwise error code 1110 will be returned.
// @Param warehouseID formData string false "1"
// @Description  type - Describes the way promotion is applied.
// @Param type formData string false "1"
// @Param awardedProductGroupID formData string false "1"
// @Param awardedBrandID formData string false "1"
// @Param storeRegionIDs formData string false "1,2,3"
// @Param customerGroupIDs formData string false "1,2,3"
// @Param lowestPriceItemIsAwarded formData string false "1 or 0"
// @Param percentageOFF formData string false "1"
// @Param sumOFF formData string false "1"
// @Param discountForOneLine formData string false "1"
// @Param requiredCouponID formData string false "1"
// @Param requiredCouponCode formData string false "1"
// @Param purchasedProducts formData string false "1"
// @Description  awardedProducts - A comma-separated list of awarded products.
// @Param awardedProducts formData string false "1,2,3,4"
// @Description  excludedProducts - A comma-separated list of excluded products.
// @Param excludedProducts formData string false "1,2,3,4"
// @Description  percentageOffExcludedProducts - A comma-separated list of percentage excluded products.
// @Param percentageOffExcludedProducts formData string false "1,2,3,4"
// @Param percentageOffIncludedProducts formData string false "1,2,3,4"
// @Param sumOffExcludedProducts formData string false "1"
// @Param sumOffIncludedProducts formData string false "1"
// @Description  awardedAmount - 	In promotion "% or $ off of specific products", how many items should get the discount. Fulfilling the promotion conditions may entitle the customer to one discounted item (awardedAmount = 1), or at most N discounted items (awardedAmount > 1), or an unlimited number of items (awardedAmount = 0). The "unlimited" option may be used in promotions such as "First item costs $3, subsequent ones are $2 each".
// @Param awardedAmount formData string false "1"
// @Param purchasedProductCategoryID formData string false "1"
// @Param awardedProductCategoryID formData string false "1"
// @Description  maximumPointsDiscount - This setting only applies to promotions that look like "Get $1 of discount for 1 loyalty point". This setting makes sure that regardless of the number of points the customer has, the points can only be exchanged for a limited amount of discount (a specified % of invoice total).
// @Param maximumPointsDiscount formData string false "1"
// @Param customerCanUseOnlyOnce formData string false "1"
// @Description  priceAtLeast - Optional, the customer must buy a certain number of items with item price more or equal to this value, doesnt work with total value or reward points. Positive decimal or 0.
// @Param priceAtLeast formData string false "1"
// @Description  priceAtMost - Optional, the customer must buy a certain number of items with item price less or equal to this value, doesnt work with total value or reward points. Positive decimal or 0.
// @Param priceAtMost formData string false "1"
// @Description  requiresManagerOverride - Set to 1 if this is a manual promotion and it should be applied to a sale with store manager's approval only. (If you attempt to set this flag on an automatic or coupon-activated promotion, error 1076 will be returned.)
// @Param requiresManagerOverride formData string false "1"
// @Param sumOffMatchingItems formData string false "1"
// @Param percentageOffMatchingItems formData string false "1"
// @Description  excludeDiscountedFromPercentageOffEntirePurchase - Indicates that the promotion should not apply to items that have already received any discount from a price list, a manual discount by the cashier, or a discount from any preceding promotion (both item-level and invoice-level promotions).
// @Param excludeDiscountedFromPercentageOffEntirePurchase formData string false "1"
// @Param excludePromotionItemsFromPercentageOffEntirePurchase formData string false "1"
// @Param awardedProductSubsidies formData string false "1"
// @Description  reasonID - Reason Code ID. A reason code can be associated with a promotion only when its purpose has been set to "PROMOTION".
// @Param reasonID formData string false "1"
// @Description  specialUnitPrice - New unit price. Field "specialUnitPrice" must not be specified together with any other award and this field is only allowed together with "purchasedAmount".
// @Param specialUnitPrice formData string false "1"
// @Description  maxItemsWithSpecialUnitPrice - Maximum limit how many items can be purchased with this special unit price. Field "maxItemsWithSpecialUnitPrice", if specified, must be equal to or larger than "purchasedAmount".
// @Param maxItemsWithSpecialUnitPrice formData string false "1"
// @Description  redemptionLimit - Maximum limit how many times the promotion can be applied to one sale. Field "redemptionLimit" is not allowed for promotions that give % off entire invoice, require reward points or apply to an unlimited number of items. Field "redemptionLimit" can only be used together with "maxItemsWithSpecialUnitPrice" (for special unit price promotions).
// @Param redemptionLimit formData string false "1"
// @Param storeGroup formData string false "1"
// @Param canBeAppliedManuallyMultipleTimes formData string false "1"
// @Param purchasedProductGroupID formData string false "1"
// @Description purchasedProductSubsidies - A comma-separated list of subsidy amounts (in euros/dollars) for each of the products specified above.
// @Param purchasedProductSubsidies formData string false "1,2,3"
// @Param purchasedBrandID formData string false "1"
// @Param purchasedAmount formData string false "1"
// @Param purchaseTotalValue formData string false "1"
// @Description  rewardPoints - The customer must redeem a specified number of loyalty points
// @Param rewardPoints formData string false "1"
// @Description  percentageOffEntirePurchase - This promotion gives a percentage discount on the entire sale.
// @Param percentageOffEntirePurchase formData string false "1"
// @Param sumOffEntirePurchase formData string false "1"
// @Param specialPrice formData string false "1"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /saveCampaign [POST]
func (saveCampaigns *SaveCampaigns) Handle(context root.IGinContext) (*response.Data, error) {

	userEntity, err := saveCampaigns.UserRepository.GetUserBySessionKey(context.PostForm("sessionKey"))
	if err != nil || userEntity.ID == 0 {
		return nil, errors.New("userEntity not found")
	}

	record, err := getRecord(context)
	if err != nil {
		return nil, err
	}

	err = saveCampaigns.CampaignHelper.Validate(record)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	c := campaign.Campaign{
		record.CampaignID,
		record.StartDate,
		record.EndDate,
		record.Name,
		record.WarehouseID,
		record.PurchasedAmount,
		record.PurchasedProductGroupID,
		record.PurchaseTotalValue,
		record.LowestPriceItemIsAwarded,
		record.SpecialPrice,
		record.PercentageOFF,
		record.SumOFF,
		record.AwardedProductGroupID,
		record.PercentageOffEntirePurchase,
		record.SumOffEntirePurchase,
		record.RewardPoints,
		record.PercentageOffMatchingItems,
		record.Type,
		now,
		userEntity.ShortName,
		0,
		"",
	}
	err = saveCampaigns.CampaignRepository.SaveCampaigns(&c)

	if err != nil {
		return nil, err
	}

	ipas, err := getInputParametersAttrs(record)
	if err != nil {
		return nil, err
	}

	attrs, err := getAttributes(ipas, &c)
	if err != nil {
		return nil, err
	}
	err = saveCampaigns.AttrsRepository.SaveAttributes(attrs)

	if err != nil {
		return nil, err
	}

	var totalRecordsCount = 0
	var recordsCount = 0
	totalRecordsCount, err = saveCampaigns.CampaignRepository.GetCampaignsCount(
		0,
		"",
	)
	if err != nil {
		return nil, err
	}
	output, err := saveCampaigns.CampaignHelper.MapToArray(append([]campaign.Campaign{}, c), map[int][]*attributes.Attribute{c.ID: attrs})
	if err != nil {
		return nil, err
	}
	recordsCount = len(output)

	return &response.Data{
		Total:           totalRecordsCount,
		TotalInResponse: recordsCount,
		Records:         output,
	}, nil
}

// New return configured struct
func New(configuration *config.Configuration) root.IRoot {

	return &SaveCampaigns{
		CampaignRepository: campaign.New(configuration),
		AttrsRepository:    attributes.New(configuration),
		CampaignHelper:     campaignhelper.New(configuration),
		UserRepository:     user.New(configuration),
		Configuration:      configuration,
	}
}

func getRecord(c root.IGinContext) (*campaignhelper.Record, error) {
	rec := campaignhelper.Record{}

	v := reflect.ValueOf(&rec)
	v = reflect.Indirect(v)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		formVal := c.PostForm(field)
		if len(formVal) == 0 {
			continue
		}
		val := v.Field(i)
		switch val.Kind() {
		case reflect.String:
			val.SetString(formVal)
		case reflect.Slice:
			s := strings.Split(formVal, ",")
			if len(s) == 0 {
				return nil, errorcodes.New(field, 1014)
			}
			slice := reflect.MakeSlice(val.Type(), 0, 0)
			x := reflect.New(slice.Type())
			x.Elem().Set(slice)
			for _, el := range s {
				switch val.Type().String() {
				case "[]int":
					v, err := strconv.Atoi(el)
					if err != nil || v < 0 {
						return nil, errorcodes.New(field, 1014)
					}
					slice = reflect.Append(slice, reflect.ValueOf(v))
				case "[]string":
					slice = reflect.Append(slice, reflect.ValueOf(el))
				default:
					log.Error(errors.New("wrong slice type: " + val.Type().String()))
				}
			}

			val.Set(slice)
		case reflect.Bool:
			if formVal == "0" || formVal == "1" {
				val.SetBool(!(formVal == "0"))
			} else {
				return nil, errorcodes.New(field, 1014)
			}
		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(formVal, 64); err != nil || f < 0 {
				return nil, errorcodes.New(field, 1014)
			} else {
				val.SetFloat(f)
			}
		case reflect.Int:
			if f, err := strconv.Atoi(formVal); err != nil || f < 0 {
				return nil, errorcodes.New(field, 1014)
			} else {
				val.SetInt(int64(f))
			}
		case reflect.Struct:
			if val.Type().String() == "time.Time" {
				if t, err := time.Parse("2006-01-02", formVal); err != nil || t.IsZero() {
					return nil, errorcodes.New(field, 1014)
				} else {
					val.Set(reflect.ValueOf(t))
				}
			}
		default:
			log.Error("wrong field type")
		}
	}
	return &rec, nil
}

func getInputParametersAttrs(r *campaignhelper.Record) (inputParametersAttrs, error) {
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
	ipa := inputParametersAttrs{}
	v = reflect.ValueOf(&ipa)
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
	return ipa, nil
}

func getAttributes(ipa inputParametersAttrs, c *campaign.Campaign) ([]*attributes.Attribute, error) {
	attrs := make([]*attributes.Attribute, 0)

	v := reflect.ValueOf(ipa)
	t := v.Type()
	v = reflect.Indirect(v)
	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		val := v.Field(i)

		if val.IsZero() {
			continue
		}
		switch val.Kind() {
		case reflect.String:
			attrs = append(attrs, &attributes.Attribute{
				ObjID:     c.ID,
				ObjTable:  "campaign",
				Name:      field,
				Type:      attributes.TEXT,
				ValueText: val.String(),
			})
		case reflect.Float32, reflect.Float64:
			attrs = append(attrs, &attributes.Attribute{
				ObjID:       c.ID,
				ObjTable:    "campaign",
				Name:        field,
				Type:        attributes.DOUBLE,
				ValueDouble: val.Float(),
			})
		case reflect.Int, reflect.Int64:
			attrs = append(attrs, &attributes.Attribute{
				ObjID:    c.ID,
				ObjTable: "campaign",
				Name:     field,
				Type:     attributes.INT,
				ValueInt: int(val.Int()),
			})
		default:
			logrus.Error(t.Field(i).Name, ": ", v.Field(i).Kind())
		}

	}

	return attrs, nil
}

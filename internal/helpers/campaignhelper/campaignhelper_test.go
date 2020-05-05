package campaignhelper

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/repositories/config"
	configMocks "github.com/zdarovich/promotion-api/internal/repositories/config/mocks"
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestCampaignHelper_Validate_NoError(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 66
	r.WarehouseID = 1

	err := ch.Validate(r)

	require.Nil(t, err)
}

func TestCampaignHelper_Validate_IsTypeValid(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.Type = "test"

	err := ch.Validate(r)
	assert.Equal(t, err, errorcodes.New("type", 1014))

	r.Type = "auto"

	err = ch.Validate(r)
	assert.Equal(t, err, nil)

	r.Type = "coupon"

	err = ch.Validate(r)
	assert.Equal(t, err, nil)

	r.Type = "manual"

	err = ch.Validate(r)
	assert.Equal(t, err, nil)
}

func TestCampaignHelper_Validate_IsStoreRegionIDsEnabled(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.StoreRegionIDs = []int{1, 2, 3}

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1028"))
}

func TestCampaignHelper_Validate_IsCustomerGroupIDsEnabled(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.CustomerGroupIDs = []int{1, 2, 3}

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1028"))
}

func TestCampaignHelper_Validate_IsStartDateValid(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.Type = "auto"
	r.PurchasedAmount = 12
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.WarehouseID = 1

	err := ch.Validate(r)
	assert.Equal(t, err, errorcodes.New("startDate", 1014))

	r.StartDate = time.Now().Add(-1 * time.Hour)

	err = ch.Validate(r)
	assert.Equal(t, err, errorcodes.New("startDate", 1014))
}

func TestCampaignHelper_Validate_IsEndDateValid(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.Type = "auto"
	r.PurchasedAmount = 12
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.WarehouseID = 1

	err := ch.Validate(r)
	assert.Equal(t, err, errorcodes.New("endDate", 1014))

	r.EndDate = time.Now().Add(-1 * time.Hour)

	err = ch.Validate(r)
	assert.Equal(t, err, errorcodes.New("endDate", 1014))
}

func TestCampaignHelper_Validate_IsRequiresManagerOverrideAndNotAutomaticOrNotCoupon(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.Type = "auto"
	r.RequiresManagerOverride = true
	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1076"))

	r.Type = "coupon"
	r.RequiresManagerOverride = true

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1076"))
}

func TestCampaignHelper_Validate_IsMultipleSetting(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 66

	r.WarehouseID = 1
	r.StoreGroup = "test"

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1110"))
}

func TestCampaignHelper_Validate_IsPurchasedProductGroupIDOrPurchasedProductCategoryIDOrPurchasedProductsAndPurchasedAmount(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1

	r.PurchasedProductGroupID = 0
	r.PurchasedProductCategoryID = 0
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1111"))

	r.PurchasedProductGroupID = 1
	r.PurchasedProductCategoryID = 0
	r.PurchasedProducts = []string{""}
	r.PurchasedAmount = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1111"))

	r.PurchasedProductGroupID = 0
	r.PurchasedProductCategoryID = 1
	r.PurchasedProducts = []string{""}
	r.PurchasedAmount = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1111"))

}

func TestCampaignHelper_Validate_IsMultipleProductOptions(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1

	r.PurchasedProductGroupID = 0
	r.PurchasedProductCategoryID = 1
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 12

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1112"))

	r.PurchasedProductGroupID = 1
	r.PurchasedProductCategoryID = 1
	r.PurchasedProducts = []string{""}
	r.PurchasedAmount = 12

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1112"))

	r.PurchasedProductGroupID = 1
	r.PurchasedProductCategoryID = 0
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 12

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1112"))
}

func TestCampaignHelper_Validate_IsAwardedProductOptionsAndSumOffOrPercentageOff(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedAmount = 66
	r.WarehouseID = 1

	r.AwardedProductGroupID = 1
	r.AwardedProductCategoryID = 0
	r.AwardedProducts = []string{}
	r.AwardedAmount = 0
	r.SumOFF = 0
	r.PercentageOFF = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1113"))

	r.AwardedProductGroupID = 0
	r.AwardedProductCategoryID = 1
	r.AwardedProducts = []string{}
	r.AwardedAmount = 0
	r.SumOFF = 0
	r.PercentageOFF = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1113"))

	r.AwardedProductGroupID = 0
	r.AwardedProductCategoryID = 0
	r.AwardedProducts = []string{"milk", "cookie"}
	r.AwardedAmount = 0
	r.SumOFF = 0
	r.PercentageOFF = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1113"))

	r.AwardedProductGroupID = 0
	r.AwardedProductCategoryID = 0
	r.AwardedProducts = []string{}
	r.AwardedAmount = 12
	r.SumOFF = 0
	r.PercentageOFF = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1113"))
}

func TestCampaignHelper_Validate_IsMultipleAwardOptions(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.AwardedProductGroupID = 0
	r.AwardedProductCategoryID = 1
	r.AwardedProducts = []string{"milk", "cookie"}
	r.AwardedAmount = 0
	r.SumOFF = 12

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1114"))

	r.AwardedProductGroupID = 1
	r.AwardedProductCategoryID = 0
	r.AwardedProducts = []string{"milk", "cookie"}
	r.AwardedAmount = 0
	r.SumOFF = 12

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1114"))

	r.AwardedProductGroupID = 0
	r.AwardedProductCategoryID = 0
	r.AwardedProducts = []string{"milk", "cookie"}
	r.AwardedAmount = 1
	r.SumOFF = 12

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1114"))
}

func TestCampaignHelper_Validate_IsPercentageExclInclProductsAndPercentageOffEntirePurchase(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.PercentageOffExcludedProducts = []string{"12", "13"}
	r.PercentageOffIncludedProducts = []string{"12", "13"}
	r.PercentageOffEntirePurchase = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1115"))
}

func TestCampaignHelper_Validate_IsSumExclInclProductsAndSumOffEntirePurchase(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.SumOffIncludedProducts = []string{"12", "13"}
	r.SumOffExcludedProducts = []string{"12", "13"}
	r.SumOffEntirePurchase = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1116"))
}

func TestCampaignHelper_Validate_IsMaximumPointsDiscountAndRewardPointsAndSumOffEntirePurchase(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.MaximumPointsDiscount = 10
	r.RewardPoints = 10
	r.SumOffEntirePurchase = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1118"))

	r.MaximumPointsDiscount = 10
	r.RewardPoints = 0
	r.SumOffEntirePurchase = 10

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1118"))

	r.MaximumPointsDiscount = 0
	r.RewardPoints = 10
	r.SumOffEntirePurchase = 10

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1118"))
}

func TestCampaignHelper_Validate_IsLowestPriceItemIsAwardedAndSumOffOrPercentageOff(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.LowestPriceItemIsAwarded = true
	r.SumOFF = 0
	r.PercentageOFF = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1119"))
}

func TestCampaignHelper_Validate_IsExcludeDiscountedFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.ExcludeDiscountedFromPercentageOffEntirePurchase = true
	r.PercentageOffEntirePurchase = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1129"))
}

func TestCampaignHelper_Validate_IsExcludePromotionItemsFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.ExcludePromotionItemsFromPercentageOffEntirePurchase = true
	r.PercentageOffEntirePurchase = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1182"))
}

func TestCampaignHelper_Validate_IsPurchasedProductSubsidiesAndPurchasedProductsAndPercentageOffMatchingItemsOrSumOffMatchingItems(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12

	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedProductSubsidies = []string{"milk", "cookie"}
	r.PercentageOffMatchingItems = 0
	r.SumOffMatchingItems = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1132"))
}

func TestCampaignHelper_Validate_IsPurchasedProductSubsidiesLenEqualsPurchasedProductsLen(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12

	r.PurchasedProducts = []string{"milk", "cookie"}
	r.PurchasedProductSubsidies = []string{"milk", "cookie", "cake"}
	r.PercentageOffMatchingItems = 12
	r.SumOffMatchingItems = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1133"))
}

func TestCampaignHelper_Validate_IsAwardedProductSubsidiesLenEqualsAwardedProductsLen(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12

	r.AwardedProductSubsidies = []string{"milk", "cookie"}
	r.AwardedProducts = []string{"milk", "cookie", "cake"}
	r.SumOFF = 12

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1134"))
}

func TestCampaignHelper_Validate_IsMaxItemsWithSpecialUnitPriceEqualsOrBiggerPurchasedAmount(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.MaxItemsWithSpecialUnitPrice = 5

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1140"))
}

func TestCampaignHelper_Validate_IsRedemptionLimitAndNotPercentageOfEntirePurchaseAndNotRewardPoints(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.RedemptionLimit = 5
	r.PercentageOffEntirePurchase = 1
	r.RewardPoints = 1
	r.SumOffEntirePurchase = 1
	r.MaximumPointsDiscount = 1

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1144"))

	r.RedemptionLimit = 5
	r.PercentageOffEntirePurchase = 1
	r.RewardPoints = 0
	r.SumOffEntirePurchase = 0
	r.MaximumPointsDiscount = 0

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1144"))

	r.RedemptionLimit = 5
	r.PercentageOffEntirePurchase = 0
	r.RewardPoints = 1
	r.SumOffEntirePurchase = 1
	r.MaximumPointsDiscount = 1

	err = ch.Validate(r)
	assert.Equal(t, err, errors.New("1144"))
}

func TestCampaignHelper_Validate_IsRedemptionLimitAndMaxItemsWithSpecialUnitPrice(t *testing.T) {
	ch := new(CampaignHelper)
	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)
	ch.ConfigRepository = cr

	r := new(Record)
	r.StartDate = time.Now().Add(1 * time.Hour)
	r.EndDate = time.Now().Add(2 * time.Hour)
	r.Type = "auto"
	r.WarehouseID = 1
	r.PurchasedAmount = 12
	r.PurchasedProducts = []string{"milk", "cookie"}

	r.RedemptionLimit = 5
	r.MaxItemsWithSpecialUnitPrice = 0

	err := ch.Validate(r)
	assert.Equal(t, err, errors.New("1145"))
}

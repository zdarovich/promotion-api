package campaignhelper

import (
	"github.com/zdarovich/promotion-api/internal/repositories/config"
	"strings"
	"time"
)

func IsPercentageOffMatchingItemsOrSumOffEntirePurchaseAndPurchasedAmount(c *Record) bool {
	if c.SumOffEntirePurchase == 0 && c.PercentageOffEntirePurchase == 0 {
		return true
	}
	return (c.SumOffEntirePurchase > 0 || c.PercentageOffEntirePurchase > 0) && c.PurchasedAmount > 0
}

func IsExcludeDiscountedFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(c *Record) bool {
	if !c.ExcludeDiscountedFromPercentageOffEntirePurchase {
		return true
	}
	return c.ExcludeDiscountedFromPercentageOffEntirePurchase && c.PercentageOffEntirePurchase > 0
}

func IsExcludePromotionItemsFromPercentageOffEntirePurchaseAndPercentageOffEntirePurchase(c *Record) bool {
	if !c.ExcludePromotionItemsFromPercentageOffEntirePurchase {
		return true
	}
	return c.ExcludePromotionItemsFromPercentageOffEntirePurchase && c.PercentageOffEntirePurchase > 0
}

func IsSpecialPriceAndPurchasedAmount(c *Record) bool {
	if c.SpecialPrice == 0 {
		return true
	}
	return c.SpecialPrice > 0 && c.PurchasedAmount > 0
}

func IsLowestPriceItemIsAwardedAndSumOffOrPercentageOff(c *Record) bool {
	if !c.LowestPriceItemIsAwarded {
		return true
	}
	return c.LowestPriceItemIsAwarded && c.SumOFF > 0 || c.PercentageOFF > 0
}

func IsRedemptionLimitAndMaxItemsWithSpecialUnitPrice(c *Record) bool {
	if c.RedemptionLimit == 0 {
		return true
	}
	return c.RedemptionLimit > 0 && c.MaxItemsWithSpecialUnitPrice > 0
}

func IsMultipleSetting(c *Record) bool {
	i := 0
	if c.StoreGroup != "" {
		i++
	}
	if c.WarehouseID != 0 {
		i++
	}
	if len(c.StoreRegionIDs) != 0 {
		i++
	}
	return i != 1
}

func IsPurchasedProductGroupIDOrPurchasedProductCategoryIDOrPurchasedProductsAndPurchasedAmount(c *Record) bool {
	if c.PurchasedProductGroupID == 0 && c.PurchasedProductCategoryID == 0 && len(c.PurchasedProducts) == 0 {
		return true
	}
	return (c.PurchasedProductGroupID != 0 || c.PurchasedProductCategoryID != 0 || len(c.PurchasedProducts) > 0) && c.PurchasedAmount > 0
}

func IsAwardedProductOptionsAndSumOffOrPercentageOff(c *Record) bool {
	if c.AwardedProductGroupID == 0 && c.AwardedProductCategoryID == 0 && c.AwardedAmount == 0 && len(c.AwardedProducts) == 0 {
		return true
	}
	return (c.AwardedProductGroupID != 0 || c.AwardedProductCategoryID != 0 || len(c.AwardedProducts) > 0 || c.AwardedAmount > 0) && (c.SumOFF > 0 || c.PercentageOFF > 0)
}

func IsPercentageExclInclProductsAndPercentageOffEntirePurchase(c *Record) bool {
	if len(c.PercentageOffExcludedProducts) == 0 && len(c.PercentageOffIncludedProducts) == 0 {
		return true
	}
	return len(c.PercentageOffExcludedProducts) > 0 && len(c.PercentageOffIncludedProducts) > 0 && c.PercentageOffEntirePurchase > 0
}
func IsSumExclInclProductsAndSumOffEntirePurchase(c *Record) bool {
	if len(c.SumOffExcludedProducts) == 0 && len(c.SumOffIncludedProducts) == 0 {
		return true
	}
	return len(c.SumOffExcludedProducts) > 0 && len(c.SumOffIncludedProducts) > 0 && c.SumOffEntirePurchase > 0
}
func IsPriceAtLeastOrPriceAtMostAndPurchasedAmount(c *Record) bool {
	if c.PriceAtMost == 0 && c.PriceAtLeast == 0 {
		return true
	}
	return (c.PriceAtLeast > 0 || c.PriceAtMost > 0) && c.PurchasedAmount > 0
}
func IsMaximumPointsDiscountAndRewardPointsAndSumOffEntirePurchase(c *Record) bool {
	if c.MaximumPointsDiscount == 0 && c.RewardPoints == 0 {
		return true
	}
	return c.MaximumPointsDiscount > 0 && c.RewardPoints > 0 && c.SumOffEntirePurchase > 0
}
func IsPurchasedProductSubsidiesAndPurchasedProductsAndPercentageOffMatchingItemsOrSumOffMatchingItems(c *Record) bool {
	if len(c.PurchasedProductSubsidies) == 0 {
		return true
	}
	return len(c.PurchasedProductSubsidies) > 0 && len(c.PurchasedProducts) > 0 && (c.PercentageOffMatchingItems > 0 || c.SumOffMatchingItems > 0)
}
func IsPurchasedProductSubsidiesLenEqualsPurchasedProductsLen(c *Record) bool {
	if len(c.PurchasedProductSubsidies) == 0 {
		return true
	}
	return len(c.PurchasedProducts) == len(c.PurchasedProductSubsidies)
}
func IsAwardedProductSubsidiesLenEqualsAwardedProductsLen(c *Record) bool {
	if len(c.AwardedProductSubsidies) == 0 {
		return true
	}
	return len(c.AwardedProducts) == len(c.AwardedProductSubsidies)
}

func IsSpecialUnitPriceAndPurchasedAmount(c *Record) bool {
	if c.SpecialUnitPrice == 0 {
		return true
	}
	return c.SpecialUnitPrice > 0 && c.PurchasedAmount > 0
}

func IsMaxItemsWithSpecialUnitPriceEqualsOrBiggerPurchasedAmount(c *Record) bool {
	if c.MaxItemsWithSpecialUnitPrice == 0 {
		return true
	}
	return c.MaxItemsWithSpecialUnitPrice >= c.PurchasedAmount
}

func IsPurchasedAmountAndPurchasedProductGroupIDOrPurchasedProductCategoryIDOrPurchasedProducts(c *Record) bool {
	if c.PurchasedAmount == 0 {
		return true
	}
	return c.PurchasedAmount > 0 && (c.PurchasedProductGroupID != 0 || c.PurchasedProductCategoryID != 0 || len(c.PurchasedProducts) > 0)
}

func IsRedemptionLimitAndNotPercentageOfEntirePurchaseAndNotRewardPoints(c *Record) bool {
	if c.RedemptionLimit == 0 {
		return true
	}
	return c.RedemptionLimit > 0 && c.PercentageOffEntirePurchase == 0 && c.RewardPoints == 0
}

func IsTypeValid(c *Record) bool {
	return c.Type == "auto" || c.Type == "manual" || c.Type == "coupon"
}

func IsStoreRegionIDsEnabled(r *Record, conf config.Conf) bool {
	if len(r.StoreRegionIDs) == 0 {
		return true
	}
	return len(r.StoreRegionIDs) != 0 &&
		strings.Contains(conf.Value, "store_regions") &&
		strings.Contains(conf.Value, "promotion_regions")
}

func IsCustomerGroupIDsEnabled(r *Record, conf config.Conf) bool {
	if len(r.CustomerGroupIDs) == 0 {
		return true
	}
	return len(r.CustomerGroupIDs) != 0 &&
		strings.Contains(conf.Value, "promotion_regions")
}

func IsStartDateValid(c *Record) bool {
	return !c.StartDate.IsZero() && c.StartDate.After(time.Now().Add(-1*time.Hour))
}

func IsEndDateValid(c *Record) bool {
	return !c.EndDate.IsZero() && c.EndDate.After(c.StartDate)
}

func IsRequiresManagerOverrideAndNotAutomaticOrNotCoupon(c *Record) bool {
	return c.RequiresManagerOverride && (c.Type == "auto" || c.Type == "coupon")
}

func IsMultipleProductOptions(c *Record) bool {
	i := 0
	if c.PurchasedProductGroupID != 0 {
		i++
	}
	if c.PurchasedProductCategoryID != 0 {
		i++
	}
	if len(c.PurchasedProducts) != 0 {
		i++
	}
	return i > 1
}
func IsMultipleAwardOptions(c *Record) bool {
	i := 0
	if c.AwardedProductGroupID != 0 {
		i++
	}
	if c.AwardedProductCategoryID != 0 {
		i++
	}
	if len(c.AwardedProducts) != 0 {
		i++
	}
	if c.AwardedAmount > 0 {
		i++
	}
	return i > 1
}
func IsReasonIdEqualsPromotion(c *Record) bool {
	// TODO
	return true
}

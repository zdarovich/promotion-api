package campaignvalidator

import "github.com/zdarovich/promotion-api/internal/repositories/campaign"

func IsMultiBuyRequirement(c *campaign.Campaign) bool {
	return !(c.PurchasedAmount <= 0 && (c.SumOffEntirePurchase > 0 || c.PercentageOffAnyOneLine > 0))
}

func IsSpecialPriceRequirement(c *campaign.Campaign) bool {
	return !(c.PurchasedAmount <= 0 && c.SpecialPrice > 0)
}

func IsLowestPriceItem(c *campaign.Campaign) bool {
	return !(c.SumOff <= 0 && c.PercentageOff <= 0 && c.AwardLowestPricedItem != 0)
}

func IsOnlyOneRequirement(c *campaign.Campaign) bool {
	var r int
	if c.PurchaseTotalValue > 0 {
		r++
	} else if c.PurchasedAmount > 0 {
		r++
	} else if c.Rewardpoints > 0 {
		r++
	}
	return r == 1
}

func IsOnlyOneAwardIsDefined(c *campaign.Campaign) bool {
	var awards int
	if c.PercentageOffAllItems > 0 {
		awards++
	} else if c.SumOffEntirePurchase > 0 {
		awards++
	} else if c.SpecialPrice > 0 {
		awards++
	} else if c.SumOffEntirePurchase > 0 {
		awards++
	} else if c.SumOff > 0 {
		awards++
	} else if c.PercentageOff > 0 {
		awards++
	} else if c.SumOffEntirePurchase > 0 {
		awards++
	}
	return awards == 1
}

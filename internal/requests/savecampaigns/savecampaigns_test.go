package savecampaigns

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"
	ctxMocks "github.com/zdarovich/promotion-api/internal/api/requests/root/mocks"
	"github.com/zdarovich/promotion-api/internal/api/response"
	config2 "github.com/zdarovich/promotion-api/internal/config"
	sqlx2 "github.com/zdarovich/promotion-api/internal/database/sqlx"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignhelper"
	"github.com/zdarovich/promotion-api/internal/repositories/attributes"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	"github.com/zdarovich/promotion-api/internal/repositories/config"
	configMocks "github.com/zdarovich/promotion-api/internal/repositories/config/mocks"
	"github.com/zdarovich/promotion-api/internal/repositories/user"
	userMocks "github.com/zdarovich/promotion-api/internal/repositories/user/mocks"
	"testing"
	"time"
)

func TestSaveCampaigns_Handle_WithNoAttributes_ReturnSuccess(t *testing.T) {
	sc := new(SaveCampaigns)
	mockDB, err := sql.Open("ramsql", "TestPromotion")
	if err != nil {
		t.Error(err)
	}
	batch := []string{
		"CREATE TABLE IF NOT EXISTS `campaign` (`id` i PRIMARY KEY, `start_date` date NOT NULL, `end_date` date NOT NULL, `name` varchar(255) NOT NULL, `warehouse_id` int(11) NOT NULL, `purchased_amount` int(11) NOT NULL, `purchased_prodgroup_id` int(11) NOT NULL, `purchase_total_value` decimal(15,2) NOT NULL, `award_lowest_priced_item` tinyint(1) NOT NULL, `special_price` decimal(15,2) NOT NULL, `percentage_off` int(11) NOT NULL, `sum_off` decimal(15,2) NOT NULL, `awarded_prodgroup_id` int(11) NOT NULL, `percentage_off_all_items` int(11) NOT NULL, `sum_off_entire_purchase` decimal(15,2) NOT NULL, `rewardpoints` int(11) NOT NULL, `percentage_off_any_one_line` int(11) NOT NULL, `type` varchar(6) NOT NULL, `added` int(11) NOT NULL, `addedby` varchar(16) NOT NULL, `changed` int(11) NOT NULL, `changedby` varchar(16) NOT NULL, PRIMARY KEY (`id`) ) ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		"CREATE TABLE IF NOT EXISTS `attributes` (`id` int(11) unsigned NOT NULL PRIMARY KEY, `obj_id` int(11) NOT NULL, `obj_table` varchar(35) NOT NULL, `name` varchar(50) NOT NULL, `type` enum('int','text','double') NOT NULL DEFAULT 'text', `value_text` varchar(255) NOT NULL, `value_int` int(11) NOT NULL, `value_double` double NOT NULL, PRIMARY KEY (`id`), KEY `obj_id` (`obj_id`), KEY `obj_table` (`obj_table`), KEY `name` (`name`), KEY `value_text` (`value_text`), KEY `value_int` (`value_int`) ) ENGINE=InnoDB DEFAULT CHARSET=latin1;",
	}
	for _, b := range batch {
		_, err = mockDB.Exec(b)
		if err != nil {
			t.Error(err)
		}
	}
	sqlxDB := sqlx.NewDb(mockDB, "ramsql")
	db := sqlx2.Mysql{
		DB: sqlxDB,
	}

	cr := new(configMocks.IRepository)
	cr.On("GetConfigByName", "vertical").Return(config.Conf{}, nil)

	ch := new(campaignhelper.CampaignHelper)
	ch.ConfigRepository = cr
	sc.CampaignHelper = ch

	c := new(config2.Configuration)
	sc.Configuration = c

	ar := new(attributes.Repository)
	sc.AttrsRepository = ar

	startDate := time.Date(2099, time.April, 12, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2099, time.April, 13, 0, 0, 0, 0, time.UTC)

	cm := new(campaign.Repository)
	cm.Database = &db
	sc.CampaignRepository = cm

	ur := new(userMocks.IRepository)
	ur.On("GetUserBySessionKey", "test").Return(user.User{ID: 1, Name: "test"}, nil)
	sc.UserRepository = ur

	ginCtx := new(ctxMocks.IGinContext)
	ginCtx.On("PostForm", "sessionKey").Return("test", nil)
	ginCtx.On("PostForm", "campaignID").Return("", nil)
	ginCtx.On("PostForm", "startDate").Return(startDate.Format("2006-01-02"), nil)
	ginCtx.On("PostForm", "endDate").Return(endDate.Format("2006-01-02"), nil)
	ginCtx.On("PostForm", "name").Return("test", nil)
	ginCtx.On("PostForm", "type").Return("auto", nil)
	ginCtx.On("PostForm", "warehouseID").Return("1", nil)

	ginCtx.On("PostForm", "awardedProductGroupID").Return("", nil)
	ginCtx.On("PostForm", "awardedBrandID").Return("", nil)
	ginCtx.On("PostForm", "lowestPriceItemIsAwarded").Return("", nil)
	ginCtx.On("PostForm", "percentageOFF").Return("", nil)
	ginCtx.On("PostForm", "sumOFF").Return("", nil)
	ginCtx.On("PostForm", "discountForOneLine").Return("", nil)
	ginCtx.On("PostForm", "requiredCouponID").Return("", nil)
	ginCtx.On("PostForm", "requiredCouponCode").Return("", nil)

	ginCtx.On("PostForm", "purchasedProducts").Return("milk,cookie", nil)

	ginCtx.On("PostForm", "awardedProducts").Return("", nil)
	ginCtx.On("PostForm", "excludedProducts").Return("", nil)
	ginCtx.On("PostForm", "percentageOffExcludedProducts").Return("", nil)
	ginCtx.On("PostForm", "percentageOffIncludedProducts").Return("", nil)
	ginCtx.On("PostForm", "purchasedProductSubsidies").Return("", nil)
	ginCtx.On("PostForm", "sumOffExcludedProducts").Return("", nil)
	ginCtx.On("PostForm", "sumOffIncludedProducts").Return("", nil)
	ginCtx.On("PostForm", "awardedProductSubsidies").Return("", nil)
	ginCtx.On("PostForm", "storeRegionIDs").Return("", nil)
	ginCtx.On("PostForm", "customerGroupIDs").Return("", nil)
	ginCtx.On("PostForm", "awardedAmount").Return("", nil)
	ginCtx.On("PostForm", "purchasedProductCategoryID").Return("", nil)
	ginCtx.On("PostForm", "awardedProductCategoryID").Return("", nil)
	ginCtx.On("PostForm", "maximumPointsDiscount").Return("", nil)
	ginCtx.On("PostForm", "customerCanUseOnlyOnce").Return("", nil)
	ginCtx.On("PostForm", "priceAtLeast").Return("", nil)
	ginCtx.On("PostForm", "priceAtMost").Return("", nil)
	ginCtx.On("PostForm", "requiresManagerOverride").Return("", nil)
	ginCtx.On("PostForm", "sumOffMatchingItems").Return("", nil)
	ginCtx.On("PostForm", "percentageOffMatchingItems").Return("", nil)
	ginCtx.On("PostForm", "excludeDiscountedFromPercentageOffEntirePurchase").Return("", nil)
	ginCtx.On("PostForm", "excludePromotionItemsFromPercentageOffEntirePurchase").Return("", nil)
	ginCtx.On("PostForm", "reasonID").Return("", nil)
	ginCtx.On("PostForm", "specialUnitPrice").Return("", nil)
	ginCtx.On("PostForm", "maxItemsWithSpecialUnitPrice").Return("", nil)
	ginCtx.On("PostForm", "redemptionLimit").Return("", nil)
	ginCtx.On("PostForm", "storeGroup").Return("", nil)
	ginCtx.On("PostForm", "canBeAppliedManuallyMultipleTimes").Return("", nil)
	ginCtx.On("PostForm", "purchasedProductGroupID").Return("", nil)
	ginCtx.On("PostForm", "purchasedBrandID").Return("", nil)

	ginCtx.On("PostForm", "purchasedAmount").Return("66", nil)

	ginCtx.On("PostForm", "purchaseTotalValue").Return("", nil)
	ginCtx.On("PostForm", "rewardPoints").Return("", nil)
	ginCtx.On("PostForm", "percentageOffEntirePurchase").Return("", nil)
	ginCtx.On("PostForm", "sumOffEntirePurchase").Return("", nil)
	ginCtx.On("PostForm", "specialPrice").Return("", nil)
	ginCtx.On("PostForm", "added").Return("", nil)
	ginCtx.On("PostForm", "addedby").Return("", nil)
	ginCtx.On("PostForm", "changed").Return("", nil)
	ginCtx.On("PostForm", "changedby").Return("", nil)

	actual, err := sc.Handle(ginCtx)
	assert.Nil(t, err)

	expected := response.Data{
		Total:           1,
		TotalInResponse: 1,
		Records: []campaignhelper.RecordOutput{
			campaignhelper.RecordOutput{
				CampaignID:        0,
				StartDate:         startDate,
				EndDate:           endDate,
				Type:              "auto",
				Name:              "test",
				WarehouseID:       1,
				PurchasedProducts: "milk,cookie",
				PurchasedAmount:   66,
				Added:             time.Now().Unix(),
			},
		},
	}

	assert.Equal(t, &expected, actual)
}

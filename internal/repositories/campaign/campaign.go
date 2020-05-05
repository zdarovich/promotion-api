package campaign

import (
	"fmt"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/database/sqlx"
	"reflect"
	"strings"
	"time"
)

type (
	// Repository struct
	Repository struct {
		Configuration *config.Configuration
		Database      sqlx.IDB
	}
	// IRepository interface
	IRepository interface {
		GetCampaigns(
			campaignID int,
			campaignType string,
			records int,
			page int,
		) ([]Campaign, error)
		GetCampaignsCount(
			campaignID int,
			campaignType string,
		) (int, error)
		SaveCampaigns(
			c *Campaign,
		) error
		UpdateCampaigns(
			c Campaign,
		) error
		DeleteCampaigns(
			campaignID int,
		) error
	}
	// Promotion structure of the promotion
	Campaign struct {
		ID                      int       `json:"id"`
		StartDate               time.Time `json:"start_date"`
		EndDate                 time.Time `json:"end_date"`
		Name                    string    `json:"name"`
		WarehouseID             int       `json:"warehouse_id"`
		PurchasedAmount         int       `json:"purchased_amount"`
		PurchasedProdgroupID    int       `json:"purchased_prodgroup_id"`
		PurchaseTotalValue      float64   `json:"purchase_total_value"`
		AwardLowestPricedItem   bool      `json:"award_lowest_priced_item"`
		SpecialPrice            float64   `json:"special_price"`
		PercentageOff           float64   `json:"percentage_off"`
		SumOff                  float64   `json:"sum_off"`
		AwardedProdgroupID      int       `json:"awarded_prodgroup_id"`
		PercentageOffAllItems   int       `json:"percentage_off_all_items"`
		SumOffEntirePurchase    float64   `json:"sum_off_entire_purchase"`
		Rewardpoints            int       `json:"rewardpoints"`
		PercentageOffAnyOneLine int       `json:"percentage_off_any_one_line"`
		Type                    string    `json:"type"`
		Added                   int64     `json:"added"`
		Addedby                 string    `json:"addedby"`
		Changed                 int64     `json:"changed"`
		Changedby               string    `json:"changedby"`
	}
)

func (repository *Repository) DeleteCampaigns(
	campaignID int,
) error {

	var query = "DELETE FROM campaign WHERE id IN (:id)"

	_, err := repository.Database.NamedExec(query,
		map[string]interface{}{
			"id": campaignID,
		})

	return err
}

func (repository *Repository) UpdateCampaigns(
	c Campaign,
) error {
	panic("implement me")
}

func (repository *Repository) SaveCampaigns(
	c *Campaign,
) error {
	conditionsString, values := repository.getSaveConditions(*c)

	var query = "INSERT INTO campaign" + conditionsString

	r, err := repository.Database.NamedExec(query, values)

	if err != nil {
		return err
	}
	id, err := r.LastInsertId()

	c.ID = int(id)

	return nil
}

func (repository *Repository) GetCampaignsCount(
	campaignID int,
	campaignType string,
) (int, error) {
	conditionsString, values := repository.getConditions(
		campaignID,
		campaignType,
	)

	var query string
	if len(conditionsString) == 0 {
		query = "SELECT COUNT(*) FROM campaign"
	} else {
		query = "SELECT COUNT(*) FROM campaign WHERE " + conditionsString
	}
	result, err := repository.Database.QueryRowx(query, values...)

	var count int

	if err != nil {
		return 0, err
	}

	err = result.Scan(&count)
	return count, err
}

// New returns new configured campaign repository
func New(configuration *config.Configuration) IRepository {

	return &Repository{
		Configuration: configuration,
		Database:      sqlx.New(configuration),
	}
}

// GetCampaigns returns campaign
func (repository *Repository) GetCampaigns(
	campaignID int,
	campaignType string,
	records int,
	page int,
) ([]Campaign, error) {

	conditionsString, values := repository.getConditions(
		campaignID,
		campaignType,
	)

	values = append(values, records*page)
	values = append(values, records*page+records)
	var query string
	if len(conditionsString) == 0 {
		query = "SELECT * FROM campaign LIMIT ?, ?"
	} else {
		query = "SELECT * FROM campaign WHERE " + conditionsString + " LIMIT ?, ?"
	}
	result, err := repository.Database.Queryx(query, values...)

	if err != nil {
		return nil, err
	}

	campaigns := make([]Campaign, 0)
	for result.Next() {
		var campaign Campaign
		err := result.StructScan(&campaign)

		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (repository *Repository) getConditions(
	campaignID int,
	campaignType string,
) (string, []interface{}) {

	var conditions []string
	var values []interface{}

	if campaignID > 0 {
		conditions = append(conditions, "id = ?")
		values = append(values, campaignID)
	}
	if campaignType != "" {
		conditions = append(conditions, "type = ?")
		values = append(values, campaignType)
	}

	return strings.Join(conditions, " AND "), values
}

func (repository *Repository) getSaveConditions(c Campaign) (string, map[string]interface{}) {
	v := reflect.ValueOf(c)
	t := v.Type()
	var fields []string
	var values []string
	vals := make(map[string]interface{})
	v = reflect.Indirect(v)
	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		if field == "id" {
			continue
		}
		fields = append(fields, field)
		values = append(values, ":"+field)
		switch v.Field(i).Kind() {
		case reflect.String:
			vals[field] = v.Field(i).String()
		case reflect.Float32, reflect.Float64:
			vals[field] = v.Field(i).Float()
		case reflect.Bool:
			vals[field] = v.Field(i).Bool()
		case reflect.Int:
			vals[field] = v.Field(i).Int()
		default:
			fmt.Println(t.Field(i).Name, ": ", v.Field(i).Kind())
			vals[field] = v.Field(i).Interface()
		}
	}
	return "(" + strings.Join(fields, ",") + ")" + " VALUES " + "(" + strings.Join(values, ",") + ")", vals
}

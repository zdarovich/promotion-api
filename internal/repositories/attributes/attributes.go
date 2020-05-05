package attributes

import (
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/database/sqlx"
	"github.com/zdarovich/promotion-api/internal/log"
)

const (
	TEXT   = "text"
	INT    = "int"
	DOUBLE = "double"
)

type (
	// Repository struct
	Repository struct {
		Configuration *config.Configuration
		Database      sqlx.IDB
	}
	// IRepository interface
	IRepository interface {
		GetAttribute(
			campaignID int,
		) ([]Attribute, error)
		GetAttributes(
			campaignIDs []int,
		) (map[int][]*Attribute, error)
		SaveAttributes(
			c []*Attribute,
		) error
		UpdateAttribute(
			c Attribute,
		) error
		DeleteAttributesByCampaignID(
			campaignID int,
		) error
	}
	// Promotion structure of the promotion
	Attribute struct {
		ID          int     `json:"id"`
		ObjID       int     `json:"obj_id"`
		ObjTable    string  `json:"obj_table"`
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		ValueText   string  `json:"value_text"`
		ValueInt    int     `json:"value_int"`
		ValueDouble float64 `json:"value_double"`
	}
)

// New returns new configured campaign repository
func New(configuration *config.Configuration) IRepository {

	return &Repository{
		Configuration: configuration,
		Database:      sqlx.New(configuration),
	}
}

func (repository *Repository) GetAttributes(
	campaignIDs []int,
) (map[int][]*Attribute, error) {
	trx, err := repository.Database.Beginx()
	defer repository.Database.Close()
	if err != nil {
		return nil, err
	}
	var query = "SELECT * FROM attributes WHERE obj_id=? AND obj_table=?"
	campaignsAttrs := make(map[int][]*Attribute)
	for _, id := range campaignIDs {
		result, err := trx.Queryx(query, id, "campaign")

		if err != nil {
			return nil, err
		}

		attrs := make([]*Attribute, 0)
		for result.Next() {
			var attr Attribute
			err := result.StructScan(&attr)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, &attr)
		}
		campaignsAttrs[id] = attrs
	}
	err = trx.Commit()
	if err != nil {
		return nil, err
	}
	return campaignsAttrs, nil
}

func (repository *Repository) GetAttribute(
	campaignID int,
) ([]Attribute, error) {

	var query = "SELECT * FROM attributes WHERE obj_id=? AND obj_table=?"

	result, err := repository.Database.Queryx(query, campaignID, "campaign")

	if err != nil {
		return nil, err
	}

	attrs := make([]Attribute, 0)
	for result.Next() {
		var attr Attribute
		err := result.StructScan(&attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, attr)
	}

	return attrs, nil
}

func (repository *Repository) SaveAttributes(
	attrs []*Attribute,
) error {
	tx, err := repository.Database.Beginx()
	defer repository.Database.Close()
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		vals := map[string]interface{}{
			"obj_id":       attr.ObjID,
			"obj_table":    attr.ObjTable,
			"name":         attr.Name,
			"type":         attr.Type,
			"value_text":   attr.ValueText,
			"value_int":    attr.ValueInt,
			"value_double": attr.ValueDouble,
		}
		res, err := tx.NamedExec("INSERT INTO attributes (obj_id, obj_table, name, type, value_text, value_int, value_double) VALUES "+
			"(:obj_id, :obj_table, :name, :type, :value_text, :value_int, :value_double)", vals)
		if err != nil {
			log.Error(tx.Rollback())
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		attr.ID = int(id)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) UpdateAttribute(
	c Attribute,
) error {
	panic("implement me")
}

func (repository *Repository) DeleteAttributesByCampaignID(
	campaignID int,
) error {
	var query = "DELETE FROM campaign WHERE obj_id=:obj_id AND obj_table=:obj_table"

	_, err := repository.Database.NamedExec(query,
		map[string]interface{}{
			"obj_id":    campaignID,
			"obj_table": "campaign",
		})

	return err
}

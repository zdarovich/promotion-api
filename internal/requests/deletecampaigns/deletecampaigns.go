package deletecampaigns

import (
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignhelper"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	"strconv"
)

type (
	// DeleteCampaigns struct
	DeleteCampaigns struct {
		CampaignRepository campaign.IRepository
		CampaignHelper     campaignhelper.ICampaignHelper
		Configuration      *config.Configuration
		InputParameters    inputParameters
	}
	// requestParams the parameters that can be used for searching
	inputParameters struct {
		CampaignID int
	}
)

// @Summary Delete campaign
// @Description  Delete campaign
// @Tags campaign
// @Accept  json
// @Produce  json
// @Param sessionKey formData string true "session key"
// @Param clientCode formData string true "client code"
// @Param request formData string true "client code"
// @Param campaignID formData string true "campaign IDs - (1,4,7)"
// @Success 200 "Created"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router / [POST]
func (deleteCampaigns *DeleteCampaigns) Handle(context root.IGinContext) (*response.Data, error) {

	err := deleteCampaigns.validate(context)

	if err != nil {
		return nil, err
	}

	var totalRecordsCount int = 0
	var recordsCount int = 0
	var records interface{}

	var campaigns []campaign.Campaign

	err = deleteCampaigns.CampaignRepository.DeleteCampaigns(
		deleteCampaigns.InputParameters.CampaignID,
	)
	if err != nil {
		return nil, err
	}
	totalRecordsCount, err = deleteCampaigns.CampaignRepository.GetCampaignsCount(
		deleteCampaigns.InputParameters.CampaignID,
		"",
	)
	if err != nil {
		return nil, err
	}
	recordsCount = len(campaigns)
	records = deleteCampaigns.CampaignHelper.MapToArray(campaigns)

	return &response.Data{
		Total:           totalRecordsCount,
		TotalInResponse: recordsCount,
		Records:         records,
	}, nil
}

// New return configured struct
func New(configuration *config.Configuration) root.IRoot {

	return &DeleteCampaigns{
		CampaignRepository: campaign.New(configuration),
		CampaignHelper:     campaignhelper.New(configuration),
		Configuration:      configuration,
	}
}

// validate checks if the required parameters have been set
func (deleteCampaigns *DeleteCampaigns) validate(context root.IGinContext) error {

	inputParameters := inputParameters{}
	inputParameters.CampaignID, _ = strconv.Atoi(context.PostForm("campaignID"))

	// Required parameters
	if inputParameters.CampaignID == 0 {

		return errors.New(strconv.Itoa(errorcodes.CodeRequiredParameterMissing))
	}

	deleteCampaigns.InputParameters = inputParameters
	return nil
}

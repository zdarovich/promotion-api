package getcampaigns

import (
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/helpers/campaignhelper"
	"github.com/zdarovich/promotion-api/internal/repositories/attributes"
	"github.com/zdarovich/promotion-api/internal/repositories/campaign"
	"strconv"
)

type (
	// GetCampaigns struct
	GetCampaigns struct {
		CampaignRepository  campaign.IRepository
		AttributeRepository attributes.IRepository
		CampaignHelper      campaignhelper.ICampaignHelper
		Configuration       *config.Configuration
		InputParameters     inputParameters
	}
	// requestParams the parameters that can be used for searching
	inputParameters struct {
		CampaignID    int
		CampaignType  string
		RecordsOnPage int
		PageNo        int
	}
)

// @Summary Get campaign
// @Description  Get campaign
// @Tags campaign
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param sessionKey formData string true "ERPLY session key"
// @Param clientCode formData string true "ERPLY client code"
// @Param request formData string true "getCampaigns"
// @Param campaignID formData string false "1"
// @Param recordsOnPage formData string false "1"
// @Param pageNo formData string false "1"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /getCampaigns [POST]
func (getCampaigns *GetCampaigns) Handle(context root.IGinContext) (*response.Data, error) {

	err := getCampaigns.validate(context)

	if err != nil {
		return nil, err
	}

	var totalRecordsCount int = 0
	var recordsCount int = 0
	var records interface{}
	campaigns, err := getCampaigns.CampaignRepository.GetCampaigns(
		getCampaigns.InputParameters.CampaignID,
		getCampaigns.InputParameters.CampaignType,
		getCampaigns.InputParameters.RecordsOnPage,
		getCampaigns.InputParameters.PageNo,
	)
	if err != nil {
		return nil, errorcodes.Wrap(err, 1003)
	}
	totalRecordsCount, err = getCampaigns.CampaignRepository.GetCampaignsCount(
		getCampaigns.InputParameters.CampaignID,
		getCampaigns.InputParameters.CampaignType,
	)
	if err != nil {
		return nil, errorcodes.Wrap(err, 1003)
	}

	attrs, err := getCampaigns.AttributeRepository.GetAttributes(campaign.GetIds(campaigns))
	if err != nil {
		return nil, errorcodes.Wrap(err, 1003)
	}
	recordsCount = len(campaigns)
	records, err = getCampaigns.CampaignHelper.MapToArray(campaigns, attrs)
	if err != nil {
		return nil, errorcodes.Wrap(err, 1003)
	}
	return &response.Data{
		Total:           totalRecordsCount,
		TotalInResponse: recordsCount,
		Records:         records,
	}, nil
}

// New return configured struct
func New(configuration *config.Configuration) root.IRoot {

	return &GetCampaigns{
		CampaignRepository:  campaign.New(configuration),
		AttributeRepository: attributes.New(configuration),
		CampaignHelper:      campaignhelper.New(configuration),
		Configuration:       configuration,
	}
}

// validate checks if the required parameters have been set
func (getCampaigns *GetCampaigns) validate(context root.IGinContext) error {

	inputParameters := inputParameters{}
	inputParameters.CampaignID, _ = strconv.Atoi(context.PostForm("campaignID"))
	inputParameters.RecordsOnPage, _ = strconv.Atoi(context.PostForm("recordsOnPage"))
	inputParameters.PageNo, _ = strconv.Atoi(context.PostForm("pageNo"))

	// Required parameters
	//if inputParameters.CampaignID == 0 {
	//
	//	return errors.New(errorcodes.CodeRequiredParameterMissing)
	//}

	// Set defaults
	if inputParameters.RecordsOnPage == 0 {
		inputParameters.RecordsOnPage = 20
	}
	if inputParameters.PageNo == 0 {
		inputParameters.PageNo = 0
	}

	getCampaigns.InputParameters = inputParameters
	return nil
}

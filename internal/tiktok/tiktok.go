package tiktok

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/mispon/tiktok-reporting-api/internal/env"
	"github.com/mispon/tiktok-reporting-api/internal/parser"
	"github.com/mispon/tiktok-reporting-api/internal/utils"
)

const (
	tokenUrl   = "https://business-api.tiktok.com/open_api/v1.2/oauth2/access_token/"
	reportsUrl = "https://business-api.tiktok.com/open_api/v1.2/reports/integrated/get/"
)

type Api interface {
	OnAuth(query url.Values) error
	GetAuctionReport(advertId int64, dateFrom string, dateTo string) (*parser.Response, error)
	GetReservationReport(advertId int64, dateFrom string, dateTo string) (*parser.Response, error)
}

type api struct {
	env    *env.Env
	parser parser.Parser
}

// New - constructor
func New(env *env.Env) Api {
	return &api{
		env:    env,
		parser: parser.New(),
	}
}

// OnAuth - process auth callback
func (a *api) OnAuth(query url.Values) error {
	authCode := query.Get("auth_code")
	fmt.Printf("[callback] received callback with auth_code: %s\n", authCode)

	if len(authCode) == 0 {
		return errors.New("received empty auth_code")
	}

	payload := map[string]interface{}{"app_id": a.env.AppId, "secret": a.env.AppSecret, "auth_code": authCode}
	jsonData, err := json.Marshal(payload)

	resp, err := utils.SendPOST(tokenUrl, jsonData)
	if err != nil {
		return err
	}

	token, advertIds, err := parseToken(resp)
	if err != nil {
		return err
	}

	a.env.AppToken = token
	a.env.AdvertiserIds = advertIds
	fmt.Printf("access token is %s\n", a.env.AppToken)

	return nil
}

// parseToken - process token response
func parseToken(respRaw map[string]interface{}) (string, []int64, error) {
	code, found := respRaw["code"]
	if !found || code.(float64) != 0 {
		return "", nil, errors.New(fmt.Sprintf("received invalid response: %v", respRaw))
	}

	data := respRaw["data"].(map[string]interface{})
	token := data["access_token"].(string)
	advertSlice := data["advertiser_ids"].([]interface{})

	advertIds := make([]int64, len(advertSlice))
	for i, ai := range advertSlice {
		advertIds[i] = int64(ai.(float64))
	}

	return token, advertIds, nil
}

// GetAuctionReport - get auction report
func (a api) GetAuctionReport(advertId int64, dateFrom string, dateTo string) (*parser.Response, error) {
	reqUrl := createUrl(advertId, dateFrom, dateTo, "AUCTION")
	respRaw, err := utils.SendGET(reqUrl, a.env.AppToken)
	if err != nil {
		return nil, err
	}

	resp, err := a.parser.Parse(respRaw)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetReservationReport - get reservation report
func (a api) GetReservationReport(advertId int64, dateFrom string, dateTo string) (*parser.Response, error) {
	reqUrl := createUrl(advertId, dateFrom, dateTo, "RESERVATION")
	respRaw, err := utils.SendGET(reqUrl, a.env.AppToken)
	if err != nil {
		return nil, err
	}

	resp, err := a.parser.Parse(respRaw)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// createUrl - creates TikTok API request URL
func createUrl(advertId int64, dateFrom string, dateTo string, serviceType string) string {
	u, _ := url.Parse(reportsUrl)

	query := u.Query()
	query.Add("service_type", serviceType)
	query.Add("report_type", "BASIC")
	query.Add("data_level", fmt.Sprintf("%s_ADVERTISER", serviceType))
	query.Add("dimensions", "[\"advertiser_id\", \"stat_time_day\"]")
	query.Add("advertiser_id", fmt.Sprintf("%d", advertId))
	query.Add("start_date", dateFrom)
	query.Add("end_date", dateTo)
	query.Add("metrics", "[\"spend\", \"ctr\", \"impressions\", \"cpc\", \"cpm\"]")
	u.RawQuery = query.Encode()

	return u.String()
}

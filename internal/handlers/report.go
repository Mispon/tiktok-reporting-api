package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mispon/tiktok-reporting-api/internal/parser"

	"github.com/mispon/tiktok-reporting-api/internal/store"

	"github.com/mispon/tiktok-reporting-api/internal/env"
	"github.com/mispon/tiktok-reporting-api/internal/utils"
)

const (
	tokenUrl   = "https://business-api.tiktok.com/open_api/v1.2/oauth2/access_token/"
	reportsUrl = "https://business-api.tiktok.com/open_api/v1.2/reports/integrated/get/"
)

type ReportHandler interface {
	Init(ctx context.Context) error
}

type reportHandler struct {
	env    *env.Env
	store  store.Store
	parser parser.Parser
	token  string
}

// New constructor
func New(env *env.Env) ReportHandler {
	return &reportHandler{
		env:    env,
		token:  env.AppToken,
		parser: parser.New(),
	}
}

// Init process handlers initialization
func (rh *reportHandler) Init(ctx context.Context) error {
	http.HandleFunc("/auth/callback", rh.callback)
	http.HandleFunc("/report/auction", rh.getAuctionReport)
	http.HandleFunc("/report/reservation", rh.getReservationReport)

	var err error
	rh.store, err = store.New(ctx, rh.env.ProjectId, rh.env.DatasetId)

	return err
}

// callback handles TikTok auth callbacks
func (rh *reportHandler) callback(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	authCode := query.Get("auth_code")
	fmt.Printf("[callback] received callback with auth_code: %s\n", authCode)

	if len(authCode) == 0 {
		io.WriteString(rw, "[callback] received empty auth_code!")
		return
	}

	payload := map[string]interface{}{"app_id": rh.env.AppId, "secret": rh.env.AppSecret, "auth_code": authCode}
	jsonData, err := json.Marshal(payload)

	resp, err := utils.SendPOST(tokenUrl, jsonData)
	if err != nil {
		fmt.Printf("[callback] failed to get token, error: %s\n", err.Error())
		return
	}

	code, found := resp["code"]
	if !found || code.(float64) != 0 {
		fmt.Printf("[callback] received invalid response: %v\n", resp)
		return
	}

	data := resp["data"].(map[string]interface{})
	rh.token = data["access_token"].(string)

	fmt.Printf("[callback] access token is %s\n", rh.token)
	io.WriteString(rw, fmt.Sprintf("%v", resp))
}

// getAuctionReport get auction marketing data from TikTok API
func (rh *reportHandler) getAuctionReport(rw http.ResponseWriter, request *http.Request) {
	requestQuery := request.URL.Query()
	fmt.Printf("[getAuctionReport] received query: %v\n", requestQuery)

	reqUrl := createUrl(&requestQuery, "AUCTION")
	respRaw, err := utils.SendGET(reqUrl, rh.token)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to get report, error: %s\n", err.Error()))
		return
	}

	resp, err := rh.parser.Parse(respRaw)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to parse API response, error: %s\n", err.Error()))
		return
	}

	err = rh.store.Save(request.Context(), resp.Data.List, rh.env.AucTableId)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to save API response, error: %s\n", err.Error()))
		return
	}

	r, err := json.Marshal(resp)
	if err == nil {
		io.WriteString(rw, string(r))
	} else {
		io.WriteString(rw, err.Error())
	}
}

// getReservationReport get reservation marketing data from TikTok API
func (rh *reportHandler) getReservationReport(rw http.ResponseWriter, request *http.Request) {
	requestQuery := request.URL.Query()
	fmt.Printf("[getReservationReport] received query: %v\n", requestQuery)

	reqUrl := createUrl(&requestQuery, "RESERVATION")
	respRaw, err := utils.SendGET(reqUrl, rh.token)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to get report, error: %s\n", err.Error()))
		return
	}

	resp, err := rh.parser.Parse(respRaw)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to parse API response, error: %s\n", err.Error()))
		return
	}

	err = rh.store.Save(request.Context(), resp.Data.List, rh.env.ResTableId)

	if err == nil {
		io.WriteString(rw, respRaw)
	} else {
		io.WriteString(rw, err.Error())
	}
}

// createUrl - creates TikTok API request URL
func createUrl(request *url.Values, serviceType string) string {
	u, _ := url.Parse(reportsUrl)

	query := u.Query()
	query.Add("service_type", serviceType)
	query.Add("report_type", "BASIC")
	query.Add("data_level", fmt.Sprintf("%s_ADVERTISER", serviceType))
	query.Add("dimensions", "[\"advertiser_id\", \"stat_time_day\"]")
	query.Add("advertiser_id", request.Get("advertiser_id"))
	query.Add("start_date", request.Get("start_date"))
	query.Add("end_date", request.Get("end_date"))
	query.Add("metrics", "[\"spend\", \"ctr\", \"impressions\", \"cpc\", \"cpm\"]")
	u.RawQuery = query.Encode()

	return u.String()
}

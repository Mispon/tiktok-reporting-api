package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	gsh "github.com/mispon/tiktok-reporting-api/internal/google_sheets"

	"github.com/mispon/tiktok-reporting-api/internal/utils"
)

const (
	tokenUrl   = "https://business-api.tiktok.com/open_api/v1.2/oauth2/access_token/"
	reportsUrl = "https://business-api.tiktok.com/open_api/v1.2/reports/integrated/get/"
)

type ReportHandler interface {
	Init()
}

// New constructor
func New(appId int, appSecret string, sbToken string) ReportHandler {
	return &reportHandler{
		AppId:        appId,
		AppSecret:    appSecret,
		Token:        sbToken,
		googleSheets: gsh.New(),
	}
}

type reportHandler struct {
	AppId        int
	AppSecret    string
	Token        string
	AdvertiserId float64
	googleSheets gsh.GoogleSheet
}

// Init process handlers initialization
func (rh *reportHandler) Init() {
	http.HandleFunc("/auth/callback", rh.callback)
	http.HandleFunc("/report/auction", rh.getAuctionReport)
	http.HandleFunc("/sheets/write_row", rh.writeRow)

	err := rh.googleSheets.Init(context.Background())
	if err != nil {
		fmt.Printf("[Init] failed to init google sheets, error: %s\n", err.Error())
	}
}

// callback handles TikTok auth callbacks
func (rh *reportHandler) callback(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	authCode := query.Get("auth_code")
	fmt.Printf("[callback] received callback with auth_code: %s\n", authCode)

	if len(authCode) == 0 {
		_, _ = io.WriteString(rw, "[callback] received empty auth_code!")
		return
	}

	payload := map[string]interface{}{"app_id": rh.AppId, "secret": rh.AppSecret, "auth_code": authCode}
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
	rh.Token = data["access_token"].(string)

	fmt.Printf("[callback] access token is %s\n", rh.Token)
	_, _ = io.WriteString(rw, fmt.Sprintf("%v", resp))
}

// getAuctionReport get auction marketing data from TikTok API
func (rh *reportHandler) getAuctionReport(rw http.ResponseWriter, request *http.Request) {
	requestQuery := request.URL.Query()
	fmt.Printf("[getAuctionReport] received query: %v\n", requestQuery)

	u, _ := url.Parse(reportsUrl)
	query := u.Query()
	query.Add("service_type", "AUCTION")
	query.Add("report_type", "BASIC")
	query.Add("data_level", "AUCTION_AD")
	query.Add("dimensions", "[\"ad_id\"]")
	query.Add("advertiser_id", requestQuery.Get("advertiser_id"))
	query.Add("start_date", requestQuery.Get("start_date"))
	query.Add("end_date", requestQuery.Get("end_date"))
	u.RawQuery = query.Encode()

	resp, err := utils.SendGET(u.String(), rh.Token)
	if err != nil {
		_, _ = io.WriteString(rw, fmt.Sprintf("failed to get report, error: %s\n", err.Error()))
		return
	}

	_, _ = io.WriteString(rw, resp)
}

// writeRow writes simple row in google spreadsheets
func (rh *reportHandler) writeRow(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	_ = request.ParseForm()

	valueRange := query.Get("value_range")
	strValues := request.Form["values"]

	var values = make([]interface{}, len(strValues))
	for i := range strValues {
		values[i] = strValues[i]
	}

	err := rh.googleSheets.WriteRow(valueRange, values)
	if err != nil {
		fmt.Printf("[writeRow] failed to write data to sheet, error: %s\n", err.Error())
	}

	_, _ = io.WriteString(rw, "ok")
}

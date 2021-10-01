package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mispon/tiktok-reporting-api/internal/utils"
)

const (
	tokenUrl = "https://business-api.tiktok.com/open_api/v1.2/oauth2/access_token/"
)

type ReportHandler interface {
	Init()
}

// New constructor
func New(appId int, appSecret string, sbToken string) ReportHandler {
	return &reportHandler{
		AppId:     appId,
		AppSecret: appSecret,
		Token:     sbToken,
	}
}

type reportHandler struct {
	AppId     int
	AppSecret string
	Token     string
}

// Init process handlers initialization
func (rh *reportHandler) Init() {
	http.HandleFunc("/auth/callback", rh.callback)
	http.HandleFunc("/report/auction", rh.getAuctionReport)
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
		fmt.Printf("[callback] failed to get token, error: %v\n", err.Error())
		return
	}

	code, found := resp["code"]
	if !found || code.(float64) != 0 {
		fmt.Printf("[callback] received invalid response: %v\n", resp)
		return
	}

	data := resp["data"].(struct{ accessToken string })
	rh.Token = data.accessToken
	fmt.Printf("[callback] access token is %s\n", rh.Token)

	_, _ = io.WriteString(rw, "ok")
}

// getAuctionReport get auction marketing data from TikTok API
func (rh *reportHandler) getAuctionReport(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fmt.Printf("[getAuctionReport] received query: %v\n", query)

	_, _ = io.WriteString(rw, "auction report data")
}

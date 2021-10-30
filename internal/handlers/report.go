package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mispon/tiktok-reporting-api/internal/tiktok"

	"github.com/mispon/tiktok-reporting-api/internal/store"

	"github.com/mispon/tiktok-reporting-api/internal/env"
)

type ReportHandler interface {
	Init(ctx context.Context, env *env.Env) error
}

type reportHandler struct {
	api        tiktok.Api
	store      store.Store
	aucTableId string
	resTableId string
}

// New constructor
func New() ReportHandler {
	return &reportHandler{}
}

// Init process handlers initialization
func (rh *reportHandler) Init(ctx context.Context, env *env.Env) error {
	var err error

	http.HandleFunc("/auth/callback", rh.callback)
	http.HandleFunc("/report/auction", rh.getAuctionReport)
	http.HandleFunc("/report/reservation", rh.getReservationReport)

	rh.aucTableId = env.AucTableId
	rh.resTableId = env.ResTableId

	rh.api = tiktok.New(env)
	rh.store, err = store.New(ctx, env.ProjectId, env.DatasetId)

	return err
}

// callback handles TikTok auth callbacks
func (rh *reportHandler) callback(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fmt.Printf("[callback] received query: %v\n", query)

	if err := rh.api.OnAuth(query); err != nil {
		io.WriteString(rw, err.Error())
	}

	io.WriteString(rw, "ok")
}

// getAuctionReport get auction marketing data from TikTok API
func (rh *reportHandler) getAuctionReport(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	advertId, dateFrom, dateTo, err := parseQuery(query)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to parse query, error: %s\n", err.Error()))
		return
	}

	resp, err := rh.api.GetAuctionReport(advertId, dateFrom, dateTo)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to get report, error: %s\n", err.Error()))
		return
	}

	err = rh.store.Save(request.Context(), resp.Data.List, rh.aucTableId)
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
	query := request.URL.Query()

	advertId, dateFrom, dateTo, err := parseQuery(query)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to parse query, error: %s\n", err.Error()))
		return
	}

	resp, err := rh.api.GetAuctionReport(advertId, dateFrom, dateTo)
	if err != nil {
		io.WriteString(rw, fmt.Sprintf("failed to get report, error: %s\n", err.Error()))
		return
	}

	err = rh.store.Save(request.Context(), resp.Data.List, rh.resTableId)
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

func parseQuery(query url.Values) (int64, string, string, error) {
	advertId, err := strconv.ParseInt(query.Get("advertiser_id"), 10, 64)
	if err != nil {
		return -1, "", "", err
	}

	dateFrom, dateTo := query.Get("start_date"), query.Get("end_date")
	if len(dateFrom) == 0 || len(dateTo) == 0 {
		return advertId, "", "", errors.New("start_date OR end_date not specified")
	}

	return advertId, dateFrom, dateTo, nil
}

package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mispon/tiktok-reporting-api/internal/parser"

	"github.com/mispon/tiktok-reporting-api/internal/store"

	"github.com/mispon/tiktok-reporting-api/internal/tiktok"

	"github.com/mispon/tiktok-reporting-api/internal/env"

	"github.com/jasonlvhit/gocron"
)

type reportGetter func() (*parser.Response, error)

type ReportJob interface {
	Schedule() error
}

type reportJob struct {
	ctx   context.Context
	env   *env.Env
	api   tiktok.Api
	store store.Store
}

func New(ctx context.Context, env *env.Env) (ReportJob, error) {
	str, err := store.New(ctx, env.ProjectId, env.DatasetId)
	if err != nil {
		return nil, err
	}

	return &reportJob{
		ctx:   ctx,
		env:   env,
		api:   tiktok.New(env),
		store: str,
	}, nil
}

// Schedule - schedule job
func (j reportJob) Schedule() error {
	s := gocron.NewScheduler()
	if err := s.Every(j.env.JobIntervalHours).Hours().Do(j.task); err != nil {
		log.Fatal(err)
	}
	<-s.Start()
	return nil
}

// task - job's work task
func (j reportJob) task() {
	if j.env.AdvertiserIds == nil {
		return
	}

	currentTime := time.Now()
	dateFrom := currentTime.AddDate(0, 0, -j.env.StatisticDepthDays).Format("2006-01-02")
	dateTo := currentTime.Format("2006-01-02")

	fmt.Printf("[%v] job started for range {from: %s, to: %s}\r\n", currentTime.Format("2006-01-02T15:04:05"), dateFrom, dateTo)

	for i, advertId := range j.env.AdvertiserIds {
		fmt.Printf("%d. processing %d advert\r\n", i+1, advertId)
		j.processStatistic(j.loadAuction(advertId, dateFrom, dateTo), j.env.AucTableId)
		j.processStatistic(j.loadReservation(advertId, dateFrom, dateTo), j.env.ResTableId)
	}

	fmt.Printf("[%v] job finished\r\n", time.Now().Format("2006-01-02T15:04:05"))
	fmt.Println()
}

// loadAuction - load auction campaigns data
func (j reportJob) loadAuction(advertId int64, dateFrom string, dateTo string) reportGetter {
	return func() (*parser.Response, error) {
		return j.api.GetAuctionReport(advertId, dateFrom, dateTo)
	}
}

// loadReservation - load reservation campaigns data
func (j reportJob) loadReservation(advertId int64, dateFrom string, dateTo string) reportGetter {
	return func() (*parser.Response, error) {
		return j.api.GetReservationReport(advertId, dateFrom, dateTo)
	}
}

// processStatistic - get and save statistic
func (j reportJob) processStatistic(getter reportGetter, tableId string) {
	response, err := getter()
	if err != nil {
		fmt.Printf("failed to get statistic, err: %v", err)
		return
	}
	fmt.Printf("\t%s -> %s\r\n", tableId, response.Message)

	if response.Data.List == nil {
		return
	}

	err = j.store.Save(j.ctx, response.Data.List, tableId)
	if err != nil {
		fmt.Printf("failed to save statistic, err: %v", err)
	}
}

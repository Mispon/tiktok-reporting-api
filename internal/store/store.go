package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/mispon/tiktok-reporting-api/internal/parser"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type Store interface {
	Save(ctx context.Context, data []parser.ListItem, tableId string) error
	Load(ctx context.Context) error
}

type store struct {
	bq        *bigquery.Client
	projectId string
	datasetId string
}

// New - constructor
func New(ctx context.Context, projectId string, datasetId string) (Store, error) {
	instance, err := bigquery.NewClient(ctx, projectId, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		return nil, err
	}
	return &store{
		bq:        instance,
		projectId: projectId,
		datasetId: datasetId,
	}, nil
}

// Save - saves data to storage
func (s *store) Save(ctx context.Context, data []parser.ListItem, tableId string) error {
	inserter := s.bq.Dataset(s.datasetId).Table(tableId).Inserter()

	items := make([]Item, len(data))
	for i, value := range data {
		items[i] = Item{
			Date:        strings.TrimSpace(strings.TrimRight(value.Dimensions.StatTimeDay, "00:00:00")),
			Spend:       value.Metrics.Spend,
			Impressions: value.Metrics.Impressions,
			Ctr:         value.Metrics.Ctr,
			Cpc:         value.Metrics.Cpc,
			Cpm:         value.Metrics.Cpm,
		}
	}

	return inserter.Put(ctx, items)
}

// Load - loads data from storage
func (s *store) Load(_ context.Context) error {
	return nil
}

func (s *store) testLoad(ctx context.Context) ([]string, error) {
	q := s.bq.Query(`
		SELECT year, name, SUM(number)
		FROM ` + "`bigquery-public-data.usa_names.usa_1910_2013`" + `
		WHERE name = "William"
		GROUP BY year, name
		ORDER BY year
	`)

	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	var values []string
	for {
		var bqValues []bigquery.Value
		if err := it.Next(&bqValues); err != nil {
			break
		}
		values = append(values, fmt.Sprintf("%v | %v | %v", bqValues[0], bqValues[1], bqValues[2]))
	}

	return values, nil
}

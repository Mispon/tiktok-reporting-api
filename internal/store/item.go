package store

import (
	"cloud.google.com/go/bigquery"
)

type Item struct {
	Date        string
	Spend       string
	Ctr         string
	Impressions string
	Cpc         string
	Cpm         string
}

func (i Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"date":  i.Date,
		"spend": i.Spend,
		"ctr":   i.Ctr,
		"imp":   i.Impressions,
		"cpc":   i.Cpc,
		"cpm":   i.Cpm,
	}, i.Date, nil
}

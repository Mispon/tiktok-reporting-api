package store

import (
	"crypto/md5"
	"fmt"

	"cloud.google.com/go/bigquery"
)

type Item struct {
	Date        string
	AdvertId    string
	Spend       string
	Ctr         string
	Impressions string
	Cpc         string
	Cpm         string
}

func (i Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"date":  i.Date,
		"ad_id": i.AdvertId,
		"spend": i.Spend,
		"ctr":   i.Ctr,
		"imp":   i.Impressions,
		"cpc":   i.Cpc,
		"cpm":   i.Cpm,
	}, i.getKey(), nil
}

func (i Item) getKey() string {
	strBytes := []byte(i.Date + i.AdvertId)
	hash := md5.Sum(strBytes)
	return fmt.Sprintf("%x", hash)
}

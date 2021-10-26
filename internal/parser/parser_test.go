package parser

import (
	"log"
	"testing"
)

const (
	responseRaw = `{
		"message": "OK",
		"code": 0,
		"data": {
			"page_info": {
				"total_number": 2,
				"page": 1,
				"page_size": 200,
				"total_page": 1
			},
			"list": [
				{
					"metrics": {
						"ad_name": "Ad name 20200923012039",
						"cpc": "116.0",
						"cpm": "52.0",
						"spend": "76.73",
						"impressions": "10505.0",
						"ctr": "1.1"
					},
					"dimensions": {
						"stat_time_day": "2020-10-17 00:00:00",
						"ad_id": 1678604629756978
					}
				},
				{
					"metrics": {
						"ad_name": "Ad name 20200923012039",
						"cpc": "132.0",
						"cpm": "34.0",
						"spend": "48.00",
						"impressions": "9805.0",
						"ctr": "1.4"
					},
					"dimensions": {
						"stat_time_day": "2020-10-16 00:00:00",
						"ad_id": 1678604629756978
					}
				}
			]
		},
		"request_id": "202011250924260101151531911200759C"
	}`
)

func TestParse(t *testing.T) {
	p := New()
	resp, err := p.Parse(responseRaw)

	if resp == nil || err != nil {
		log.Fatalf(`parser.Parse() failed, %v`, err)
	}
}

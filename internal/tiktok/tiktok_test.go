package tiktok

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tokenResponseJson = `{
		"message": "OK",
		"code": 0,
		"data": {
			"access_token": "ccdfaeef518606eb3f0e25355be92f51f13c3564",
			"scope": [10, 4, 15],
			"advertiser_ids": [
				6996974816561479681,
				6999955650176401410,
				7019189042788958210,
				7006232246877241345,
				7016495124368588801,
				7016498340581867521,
				7002523003154153474,
				7017662568659451906,
				7017663796814577666,
				7016533236813742082
			]
		},
		"request_id": "202110270708240102452480431884BB93"}
	`
)

func TestParseToken(t *testing.T) {
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(tokenResponseJson), &resp); err != nil {
		log.Fatalf("failed to unmarshall reponse, %v", err)
	}

	token, advertIds, err := parseToken(resp)

	assert.True(t, token == "ccdfaeef518606eb3f0e25355be92f51f13c3564")
	assert.True(t, len(advertIds) == 10)
	assert.Nil(t, err)
}

package env

import (
	"os"
	"strconv"
)

type Env struct {
	Endpoint   string
	AppId      int
	AppSecret  string
	AppToken   string
	ProjectId  string
	DatasetId  string
	AucTableId string
	ResTableId string
}

func New() Env {
	return Env{
		Endpoint:   getStr("API_ENDPOINT", "0.0.0.0:80"),
		AppId:      getInt("TIKTOK_APP_ID", 0),
		AppSecret:  getStr("TIKTOK_APP_SECRET", ""),
		AppToken:   getStr("TIKTOK_APP_TOKEN", ""),
		ProjectId:  getStr("BQ_PROJECT_ID", ""),
		DatasetId:  getStr("BQ_DATASET_ID", ""),
		AucTableId: getStr("BQ_AUC_TABLE_ID", ""),
		ResTableId: getStr("BQ_RES_TABLE_ID", ""),
	}
}

func getStr(key string, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	valueStr := getStr(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

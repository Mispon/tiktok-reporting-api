package parser

// Response - TikTok API response model
type Response struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Data      Data   `json:"data"`
	RequestId string `json:"request_id"`
}

type Data struct {
	PageInfo PageInfo   `json:"page_info"`
	List     []ListItem `json:"list"`
}

type PageInfo struct {
	TotalNumber int `json:"total_number"`
	Page        int `json:"page"`
	PageSize    int `json:"page_size"`
	TotalPage   int `json:"total_page"`
}

type ListItem struct {
	Metrics    Metrics    `json:"metrics"`
	Dimensions Dimensions `json:"dimensions"`
}

type Metrics struct {
	Spend       string `json:"spend"`
	Impressions string `json:"impressions"`
	Ctr         string `json:"ctr"`
	Cpc         string `json:"cpc"`
	Cpm         string `json:"cpm"`
}

type Dimensions struct {
	StatTimeDay  string `json:"stat_time_day"`
	AdvertiserId int64  `json:"advertiser_id"`
}

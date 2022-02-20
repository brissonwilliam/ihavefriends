package models

type BonneFeteAnalytics struct {
	Total        uint            `json:"total"`
	TotalByUsers []BFTotalByUser `json:"totalByUsers"`
}

type BFTotalByUser struct {
	Username string `json:"name" db:"username"`
	Count    uint   `json:"count" db:"total"`
}

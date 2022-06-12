package models

import (
	"strconv"
	"time"

	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
)

type Amount float64

func (a Amount) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(a), 'f', 2, 64)), nil
}

type CreateBill struct {
	ID           uuid.OrderedUUID `json:"-" db:"id" validate:"required"`
	UserId       uuid.OrderedUUID `json:"-" db:"user_id" validate:"required"`
	NewBillTotal float64          `json:"total" db:"new_bill_total" validate:"required"`
}

type BillAnalytics struct {
	GrandTotal   Amount               `json:"grandTotal"`
	TotalByUsers []BillUserTotal      `json:"totalByUsers"`
	UserTotal    BillUserTotal        `json:"userTotal"`
	TotalsByTime BillUserTotalsByTime `json:"userTotalsByTime"`
}

type BillUserTotal struct {
	UserId          uuid.OrderedUUID `json:"-" db:"user_id"`
	Username        string           `json:"name" db:"username"`
	CumulativeTotal Amount           `json:"total" db:"cumulative_total"`
	HighestTotal    Amount           `json:"highestTotal" db:"highest_total"`
	LastTotal       Amount           `json:"lastTotal" db:"last_total"`
	LastVisit       *time.Time       `json:"lastVisit" db:"last_updated"`
}

type BillUserTotalsByTime struct {
	UserId         uuid.OrderedUUID `json:"-" db:"user_id"`
	Last48HTotal   Amount           `json:"last48hTotal" db:"last_48h_total"`
	LastWeekTotal  Amount           `json:"lastWeekTotal" db:"last_week_total"`
	ThisWeekTotal  Amount           `json:"thisWeekTotal" db:"this_week_total"`
	ThisMonthTotal Amount           `json:"thisMonthTotal" db:"this_month_total"`
}

type MonthlyBillTotals struct {
}

type UserMonthlyBillTotals struct {
}

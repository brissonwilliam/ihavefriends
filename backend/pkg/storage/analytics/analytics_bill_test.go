package analytics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetLastMonday(t *testing.T) {
	type test struct {
		inDate      time.Time
		expectedOut time.Time
		description string
	}

	nonUtcLocation := time.FixedZone("CST", 8*3600)

	tests := []test{
		{
			inDate:      time.Date(2022, 04, 19, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "tuesday",
		},
		{
			inDate:      time.Date(2022, 04, 20, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "wednesday",
		},
		{
			inDate:      time.Date(2022, 04, 21, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "thursday",
		},
		{
			inDate:      time.Date(2022, 04, 22, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "friday",
		},
		{
			inDate:      time.Date(2022, 04, 23, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "saturday",
		},
		{
			inDate:      time.Date(2022, 04, 24, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 18, 4, 0, 0, 0, nonUtcLocation),
			description: "sunday",
		},
		{
			inDate:      time.Date(2022, 04, 25, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 25, 4, 0, 0, 0, nonUtcLocation),
			description: "monday (current date)",
		},
		{
			inDate:      time.Date(2022, 04, 01, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 03, 28, 4, 0, 0, 0, nonUtcLocation),
			description: "gets monday from another month",
		},
	}

	for _, unitTest := range tests {
		lastMonday := getStartOfWeek(unitTest.inDate)
		assert.Equal(t, unitTest.expectedOut, lastMonday, unitTest.description)
	}
}

func TestGetStartOfMonth(t *testing.T) {
	type test struct {
		inDate      time.Time
		expectedOut time.Time
		description string
	}

	nonUtcLocation := time.FixedZone("CST", 8*3600)

	tests := []test{
		{
			inDate:      time.Date(2022, 04, 19, 0, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 01, 4, 0, 0, 0, nonUtcLocation),
			description: "aligns to first day of month",
		},
	}

	for _, unitTest := range tests {
		lastMonday := getStartOfMonth(unitTest.inDate)
		assert.Equal(t, unitTest.expectedOut, lastMonday, unitTest.description)
	}
}

func TestGetEndOfDay(t *testing.T) {
	type test struct {
		inDate      time.Time
		expectedOut time.Time
		description string
	}

	nonUtcLocation := time.FixedZone("CST", 8*3600)

	tests := []test{
		{
			inDate:      time.Date(2022, 04, 19, 1, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 04, 19, 4, 0, 0, 0, nonUtcLocation),
			description: "keeps current day if under end of day hour",
		},
		{
			inDate:      time.Date(2022, 04, 30, 4, 0, 0, 0, nonUtcLocation),
			expectedOut: time.Date(2022, 05, 1, 4, 0, 0, 0, nonUtcLocation),
			description: "adds day if over or equal end of day hour",
		},
	}

	for _, unitTest := range tests {
		lastMonday := getEndOfDay(unitTest.inDate)
		assert.Equal(t, unitTest.expectedOut, lastMonday, unitTest.description)
	}
}

package controllers

import (
	"strconv"
	"testing"
	"time"
)

func Test_UpgradeTimeReached(t *testing.T) {
	testCases := []struct {
		name    string
		time    time.Time
		reached bool
	}{
		{
			name:    "case 0",
			time:    time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC),
			reached: false,
		},
		{
			name:    "case 1",
			time:    time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC),
			reached: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := upgradeTimeReached(tc.time)

			if result != tc.reached {
				t.Fatalf("%s -  expected '%t' got '%t'\n", tc.name, tc.reached, result)
			}
		})
	}
}

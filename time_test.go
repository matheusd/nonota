package nonota

import (
	"time"
	"testing"
)

func TestDayTimes(t *testing.T) {
	const format = "2006-01-02 15:04:05"

	type testCase struct {
		t string
		start string
		end string
	}
	testCases := []testCase{
		// Random date (gmt -03:00 at BRT)
		{"2019-03-29 10:20:23", "2019-03-29 00:00:00", "2019-03-29 23:59:59"},
		{"2019-03-01 14:20:23", "2019-03-01 00:00:00", "2019-03-01 23:59:59"},
		{"2019-03-01 00:00:00", "2019-03-01 00:00:00", "2019-03-01 23:59:59"},
		{"2019-03-01 23:59:59", "2019-03-01 00:00:00", "2019-03-01 23:59:59"},
		{"2019-03-31 00:00:00", "2019-03-31 00:00:00", "2019-03-31 23:59:59"},
		{"2019-03-31 23:59:59", "2019-03-31 00:00:00", "2019-03-31 23:59:59"},

		// End of year (gmt -02:00 at BRT)
		{"2019-12-29 10:20:23", "2019-12-29 00:00:00", "2019-12-29 23:59:59"},
		{"2019-12-01 14:20:23", "2019-12-01 00:00:00", "2019-12-01 23:59:59"},
		{"2019-12-01 00:00:00", "2019-12-01 00:00:00", "2019-12-01 23:59:59"},
		{"2019-12-01 23:59:59", "2019-12-01 00:00:00", "2019-12-01 23:59:59"},
		{"2019-12-31 00:00:00", "2019-12-31 00:00:00", "2019-12-31 23:59:59"},
		{"2019-12-31 23:59:59", "2019-12-31 00:00:00", "2019-12-31 23:59:59"},

		// Non leap year feb date
		{"2019-02-08 19:20:23", "2019-02-08 00:00:00", "2019-02-08 23:59:59"},

		// Leap year feb date
		{"2020-02-08 05:20:23", "2020-02-08 00:00:00", "2020-02-08 23:59:59"},
	}

	for i, tc := range testCases {
		testTime, err := time.ParseInLocation(format, tc.t, time.Local)
		if err != nil {
			t.Fatalf("unable to decode test time %s: %v", tc.t, err)
		}

		expectedEnd, err := time.ParseInLocation(format, tc.end, time.Local)
		if err != nil {
			t.Fatalf("unable to decode end time %s: %v", tc.end, err)
		}

		expectedStart, err := time.ParseInLocation(format, tc.start, time.Local)
		if err != nil {
			t.Fatalf("unable to decode start time %s: %v", tc.start, err)
		}

		actualEnd := EndOfDay(testTime)
		actualStart := StartOfDay(testTime)
		if actualStart != expectedStart {
			t.Fatalf("%d (tc %s start): expected %s found %s", i, tc.t, expectedStart,
				actualStart)
		}
		if actualEnd != expectedEnd {
			t.Fatalf("%d (tc %s end): expected %s found %s", i, tc.t, expectedEnd,
				actualEnd)
		}
	}
}

func TestWeekTimes(t *testing.T) {
	const format = "2006-01-02 15:04:05"

	type testCase struct {
		t string
		start string
		end string
	}
	testCases := []testCase{
		{"2019-03-23 10:20:23", "2019-03-17 00:00:00", "2019-03-23 23:59:59"},
		{"2019-03-24 10:20:23", "2019-03-24 00:00:00", "2019-03-30 23:59:59"},
		{"2019-03-26 10:20:23", "2019-03-24 00:00:00", "2019-03-30 23:59:59"},
		{"2019-03-29 10:20:23", "2019-03-24 00:00:00", "2019-03-30 23:59:59"},
		{"2019-03-30 11:20:23", "2019-03-24 00:00:00", "2019-03-30 23:59:59"},
		{"2019-03-31 11:20:23", "2019-03-31 00:00:00", "2019-04-06 23:59:59"},
	}

	for i, tc := range testCases {
		testTime, err := time.ParseInLocation(format, tc.t, time.Local)
		if err != nil {
			t.Fatalf("unable to decode test time %s: %v", tc.t, err)
		}

		expectedEnd, err := time.ParseInLocation(format, tc.end, time.Local)
		if err != nil {
			t.Fatalf("unable to decode end time %s: %v", tc.end, err)
		}

		expectedStart, err := time.ParseInLocation(format, tc.start, time.Local)
		if err != nil {
			t.Fatalf("unable to decode start time %s: %v", tc.start, err)
		}

		actualEnd := EndOfWeek(testTime)
		actualStart := StartOfWeek(testTime)
		if actualStart != expectedStart {
			t.Fatalf("%d (tc %s start): expected %s found %s", i, tc.t, expectedStart,
				actualStart)
		}
		if actualEnd != expectedEnd {
			t.Fatalf("%d (tc %s end): expected %s found %s", i, tc.t, expectedEnd,
				actualEnd)
		}
	}
}

func TestBillTimes(t *testing.T) {
	const format = "2006-01-02 15:04:05"

	type testCase struct {
		t string
		start string
		end string
	}
	testCases := []testCase{
		{"2019-02-23 10:20:23", "2019-02-01 00:00:00", "2019-02-28 23:59:59"},
		{"2019-03-23 10:20:23", "2019-03-01 00:00:00", "2019-03-31 23:59:59"},
		{"2019-12-23 10:20:23", "2019-12-01 00:00:00", "2019-12-31 23:59:59"},
		{"2020-02-23 10:20:23", "2020-02-01 00:00:00", "2020-02-29 23:59:59"},
	}

	for i, tc := range testCases {
		testTime, err := time.ParseInLocation(format, tc.t, time.Local)
		if err != nil {
			t.Fatalf("unable to decode test time %s: %v", tc.t, err)
		}

		expectedEnd, err := time.ParseInLocation(format, tc.end, time.Local)
		if err != nil {
			t.Fatalf("unable to decode end time %s: %v", tc.end, err)
		}

		expectedStart, err := time.ParseInLocation(format, tc.start, time.Local)
		if err != nil {
			t.Fatalf("unable to decode start time %s: %v", tc.start, err)
		}

		actualEnd := EndOfBilling(testTime)
		actualStart := StartOfBilling(testTime)
		if actualStart != expectedStart {
			t.Fatalf("%d (tc %s start): expected %s found %s", i, tc.t, expectedStart,
				actualStart)
		}
		if actualEnd != expectedEnd {
			t.Fatalf("%d (tc %s end): expected %s found %s", i, tc.t, expectedEnd,
				actualEnd)
		}
	}
}
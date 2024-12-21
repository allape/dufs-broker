package env

import "testing"

type PortRangeTestCase struct {
	PortRange PortRange
	Start     int
	End       int
	hasError  bool
}

func TestPortRange(t *testing.T) {
	var start, end int
	var err error

	portRanges := []PortRangeTestCase{
		{PortRange("50000-50100"), 50000, 50100, false},
		{PortRange("50000-50000"), 50000, 50000, true},
		{PortRange("50000-49999"), 0, 0, true},
		{PortRange("2020-"), 0, 0, true},
		{PortRange("-2021"), 0, 0, true},
		{PortRange("12345"), 0, 0, true},
		{PortRange("50000-50100-50200"), 0, 0, true},
		{PortRange(""), 0, 0, true},
		{PortRange("1231sadwd"), 0, 0, true},
	}

	for _, portRange := range portRanges {
		start, end, err = portRange.PortRange.Range()
		if portRange.hasError {
			if err == nil {
				t.Errorf("Expected error but got nil")
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if start != portRange.Start {
			t.Errorf("Expected start %d but got %d", portRange.Start, start)
		}

		if end != portRange.End {
			t.Errorf("Expected end %d but got %d", portRange.End, end)
		}
	}
}

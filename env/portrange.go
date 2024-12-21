package env

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidPortRange = errors.New("invalid port range")

type PortRange string

func (r PortRange) Range() (int, int, error) {
	segments := strings.Split(string(r), "-")
	if len(segments) != 2 {
		return 0, 0, ErrInvalidPortRange
	}

	start, err := strconv.Atoi(segments[0])
	if err != nil {
		return 0, 0, ErrInvalidPortRange
	} else if start < 0 {
		return 0, 0, ErrInvalidPortRange
	}

	end, err := strconv.Atoi(segments[1])
	if err != nil || end < 0 {
		return 0, 0, ErrInvalidPortRange
	}

	if start >= end {
		return 0, 0, ErrInvalidPortRange
	}

	return start, end, nil
}

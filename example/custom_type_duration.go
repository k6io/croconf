package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)


// ParseExtendedDuration is a helper function that allows for string duration
// values containing days.
func ParseExtendedDuration(data string) (result time.Duration, err error) {
	// Assume millisecond values if data is provided with no units
	if t, errp := strconv.ParseFloat(data, 64); errp == nil {
		return time.Duration(t * float64(time.Millisecond)), nil
	}

	dPos := strings.IndexByte(data, 'd')
	if dPos < 0 {
		return time.ParseDuration(data)
	}

	var hours time.Duration
	if dPos+1 < len(data) { // case "12d"
		hours, err = time.ParseDuration(data[dPos+1:])
		if err != nil {
			return
		}
		if hours < 0 {
			return 0, fmt.Errorf("invalid time format '%s'", data[dPos+1:])
		}
	}

	days, err := strconv.ParseInt(data[:dPos], 10, 64)
	if err != nil {
		return
	}
	if days < 0 {
		hours = -hours
	}
	return time.Duration(days)*24*time.Hour + hours, nil
}


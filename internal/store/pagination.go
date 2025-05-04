package store

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TaskFilters struct {
	Date      *string `json:"date"` // "today", "yesterday", or "YYYY-MM-DD"
	Limit     int     `json:"limit" validate:"gte=1,lte=100"`
	Offset    int     `json:"offset" validate:"gte=0"`
	Completed *bool   `json:"completed"`
}

type GoalFilters struct {
	Limit    int64 `json:"limit" validate:"gte=1,lte=100"`
	Offset   int64 `json:"offset" validate:"gte=0"`
	Achieved *bool `json:"completed"`
}

func (tf TaskFilters) Parse(r *http.Request) (TaskFilters, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return tf, err
		}
		tf.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return tf, err
		}
		tf.Offset = l
	}

	if date := qs.Get("date"); date != "" {
		switch date {
		case "today":
			today := time.Now().UTC().Format("2006-01-02")
			tf.Date = &today
		case "yesterday":
			yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")
			tf.Date = &yesterday
		default:
			if _, err := time.Parse("2006-01-02", date); err != nil {
				return tf, fmt.Errorf("invalid date format: use today/yesterday or YYYY-MM-DD")
			}
			tf.Date = &date
		}
	}

	completed := qs.Get("completed")
	if completed != "" {
		c, err := strconv.ParseBool(completed)
		if err != nil {
			return tf, err
		}
		tf.Completed = &c
	}

	return tf, nil
}

func (gf GoalFilters) Parse(r *http.Request) (GoalFilters, error) {
	qs := r.URL.Query()

	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil || l < 1 || l > 100 {
			return gf, fmt.Errorf("limit must be between 1-100")
		}
		gf.Limit = l
	}

	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.ParseInt(offset, 10, 64)
		if err != nil || o < 0 {
			return gf, fmt.Errorf("offset must be positive")
		}
		gf.Offset = o
	}

	if achieved := qs.Get("achieved"); achieved != "" {
		a, err := strconv.ParseBool(achieved)
		if err != nil {
			return gf, fmt.Errorf("invalid achieved value")
		}
		gf.Achieved = &a
	}

	return gf, nil
}

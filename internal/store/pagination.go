package store

import (
	"net/http"
	"strconv"
	"strings"
	// "time"
)

type PaginationFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	// Since  string   `json:"since"`
	// Until  string   `json:"until"`
}

func (pg PaginationFeedQuery) Parse(r *http.Request) (PaginationFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pg, nil
		}
		pg.Limit = l
	}
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return pg, nil
		}
		pg.Offset = o
	}
	sort := qs.Get("sort")
	if sort != "" {
		pg.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		pg.Tags = strings.Split(tags, ",")
	}
	// since := qs.Get("since")
	// if since != "" {
	// 	pg.Since = parseTime(since)
	// }
	// until := qs.Get("until")
	// if until != "" {
	// 	pg.Until = parseTime(until)
	// }
	return pg, nil
}

// func parseTime(s string) string {
// 	t, err := time.Parse(time.DateTime, s)
// 	if err != nil {
// 		return ""
// 	}
// 	return t.Format(time.DateTime)
// }

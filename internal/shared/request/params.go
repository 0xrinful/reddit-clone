package request

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

func ReadID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}
	return id, nil
}

func ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func ReadCursor(r *http.Request, v *validator.Validator) pagination.Cursor {
	c := pagination.Cursor{Limit: pagination.DefaultLimit}

	if s := r.URL.Query().Get("after"); s != "" {
		if after, err := strconv.ParseInt(s, 10, 64); err == nil {
			if after < 1 {
				v.AddError("after", "must be greater than zero")
			} else {
				c.After = after
			}
		} else {
			v.AddError("after", "must be an integer value")
		}
	}

	if s := r.URL.Query().Get("limit"); s != "" {
		if limit, err := strconv.Atoi(s); err == nil {
			if limit > pagination.MaxLimit {
				limit = pagination.MaxLimit
			}
			if limit < 1 {
				v.AddError("limit", "must be greater than zero")
			} else {
				c.Limit = limit
			}
		} else {
			v.AddError("limit", "must be an integer value")
		}
	}
	return c
}

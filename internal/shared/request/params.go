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

func ParseCursorPagination[T any](
	r *http.Request,
	v *validator.Validator,
	decodeFunc func(string) (*T, error),
) pagination.CursorParams[T] {
	params := pagination.CursorParams[T]{Limit: pagination.DefaultLimit}
	if s := r.URL.Query().Get("cursor"); s != "" {
		cursor, err := decodeFunc(s)
		if err != nil {
			v.AddError("cursor", "invalid cursor value")
		} else {
			params.Cursor = cursor
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
				params.Limit = limit
			}
		} else {
			v.AddError("limit", "must be an integer value")
		}
	}
	return params
}

func ParseOffsetPagination(
	r *http.Request,
	v *validator.Validator,
) pagination.OffsetParams {
	params := pagination.OffsetParams{Limit: pagination.DefaultLimit, Page: 1}
	if s := r.URL.Query().Get("page"); s != "" {
		if page, err := strconv.Atoi(s); err == nil {
			v.Check(page > 0, "page", "must be greater than zero")
			v.Check(page <= 10_000_000, "page", "must be a maximum of 10 million")
			if v.Valid() {
				params.Page = page
			}
		} else {
			v.AddError("limit", "must be an integer value")
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
				params.Limit = limit
			}
		} else {
			v.AddError("limit", "must be an integer value")
		}
	}
	return params
}

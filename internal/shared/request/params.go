package request

import (
	"net/http"
	"strconv"
)

func ReadID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}
	return id, nil
}

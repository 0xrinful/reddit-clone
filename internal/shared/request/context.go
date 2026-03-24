package request

import (
	"context"
	"net/http"
)

type (
	communityKey struct{}
	userKey      struct{}
)

type CommunityCtx struct {
	ID   int64
	Name string
}

func WithCommunity(r *http.Request, c *CommunityCtx) *http.Request {
	ctx := context.WithValue(r.Context(), communityKey{}, c)
	return r.WithContext(ctx)
}

func GetCommunity(r *http.Request) *CommunityCtx {
	community, ok := r.Context().Value(communityKey{}).(*CommunityCtx)
	if !ok {
		panic("request.GetCommunity: LoadCommunity middleware not registered on this route")
	}
	return community
}

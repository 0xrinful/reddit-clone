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
		panic("missing community_name value in request context")
	}
	return community
}

type UserCtx struct {
	ID int64
}

func WithUser(r *http.Request, u *UserCtx) *http.Request {
	ctx := context.WithValue(r.Context(), userKey{}, u)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) (*UserCtx, bool) {
	u, ok := r.Context().Value(userKey{}).(*UserCtx)
	return u, ok
}

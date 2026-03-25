package posts

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/shared/apperr"
)

type Service interface {
	GetPost(ctx context.Context, id, communityID int64) (*Post, error)
	CreatePost(ctx context.Context, userID, communityID int64, req CreatePostRequest) (*Post, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreatePost(
	ctx context.Context,
	userID, communityID int64,
	req CreatePostRequest,
) (*Post, error) {
	p := &Post{
		Title:       req.Title,
		Body:        req.Body,
		UserID:      userID,
		CommunityID: communityID,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) GetPost(ctx context.Context, id, communityID int64) (*Post, error) {
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.CommunityID != communityID {
		return nil, apperr.ErrNotFound
	}

	return post, nil
}

package posts

import (
	"context"
)

type Service interface {
	GetPost(ctx context.Context, id int64) (*Post, error)
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

func (s *service) GetPost(ctx context.Context, id int64) (*Post, error) {
	return nil, nil
}

package posts

import (
	"context"
)

type Service interface {
	Get(ctx context.Context, id int64) (*PostView, error)
	Create(ctx context.Context, params CreateParams) (*Post, error)
	Update(ctx context.Context, id, actorID int64, p UpdateParams) (*Post, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, params ListParams) ([]*PostSummary, error)
}

type Authorizer interface {
	CanPost(ctx context.Context, communityID, userID int64) error
	CanUpdatePost(ctx context.Context, communityID, actorID, authorID int64) error
	CanDeletePost(actorID, authorID int64) error
}

type service struct {
	authz     Authorizer
	postsRepo Repository
}

func NewService(authz Authorizer, postsRepo Repository) Service {
	return &service{authz: authz, postsRepo: postsRepo}
}

func (s *service) Create(ctx context.Context, params CreateParams) (*Post, error) {
	if err := s.authz.CanPost(ctx, params.CommunityID, params.UserID); err != nil {
		return nil, err
	}

	p := &Post{
		Title:       params.Title,
		Body:        params.Body,
		UserID:      params.UserID,
		CommunityID: params.CommunityID,
	}

	if err := s.postsRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) Update(ctx context.Context, id, actorID int64, p UpdateParams) (*Post, error) {
	post, err := s.postsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.authz.CanUpdatePost(ctx, post.CommunityID, actorID, post.UserID); err != nil {
		return nil, err
	}

	return s.postsRepo.Update(ctx, id, p)
}

func (s *service) Get(ctx context.Context, id int64) (*PostView, error) {
	return s.postsRepo.GetView(ctx, id)
}

func (s *service) Delete(ctx context.Context, id, actorID int64) error {
	post, err := s.postsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authz.CanDeletePost(actorID, post.UserID); err != nil {
		return err
	}

	return s.postsRepo.Delete(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	params ListParams,
) ([]*PostSummary, error) {
	return s.postsRepo.List(ctx, params)
}

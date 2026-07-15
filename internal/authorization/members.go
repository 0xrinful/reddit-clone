package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
)

func (s *service) CanViewPermissions(
	ctx context.Context,
	communityID, actorID int64,
) (bool, error) {
	actor, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return false, err
	}

	return actor.Role.AtLeast(domain.RoleModerator), nil
}

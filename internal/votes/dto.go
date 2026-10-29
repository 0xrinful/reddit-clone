package votes

import (
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type VotePostRequest struct {
	Value int8 `json:"value"`
}

func (r *VotePostRequest) Validate(v *validator.Validator) {
	v.Check(VoteValue(r.Value).Valid(), "value", "must be 1 (upvote) or -1 (downvote)")
}

type VoteCommentRequest struct {
	Value int8 `json:"value"`
}

func (r *VoteCommentRequest) Validate(v *validator.Validator) {
	v.Check(VoteValue(r.Value).Valid(), "value", "must be 1 (upvote) or -1 (downvote)")
}

// DTOs
// mapping helpers
// response envelope
type VotePostResponse struct {
	PostID int64     `json:"post_id"`
	Score  int64     `json:"score"`
	Value  VoteValue `json:"value"`
}

type VoteCommentResponse struct {
	CommentID int64     `json:"comment_id"`
	Score     int64     `json:"score"`
	Value     VoteValue `json:"value"`
}

// response constructor
func toVotePostResponse(result *VotePostResult) VotePostResponse {
	return VotePostResponse{PostID: result.PostID, Score: result.Score, Value: result.Value}
}

func toVoteCommentResponse(result *VoteCommentResult) VoteCommentResponse {
	return VoteCommentResponse{
		CommentID: result.CommentID,
		Score:     result.Score,
		Value:     result.Value,
	}
}

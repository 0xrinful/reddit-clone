package votes

type VoteValue int8

const (
	Upvote   VoteValue = 1
	Downvote VoteValue = -1
)

func (v VoteValue) Valid() bool {
	return v == Upvote || v == Downvote
}

type PostVote struct {
	UserID int64
	PostID int64
	Value  VoteValue
}

type VotePostResult struct {
	PostID int64
	Score  int64
	Value  VoteValue
}

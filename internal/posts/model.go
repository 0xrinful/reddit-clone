package posts

import "time"

type Post struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Body      string
}

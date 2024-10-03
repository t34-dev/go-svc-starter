package jsonplaceholder

import (
	"context"
)

// JSONPlaceholderService defines the interface for interacting with the JSONPlaceholder API
type JSONPlaceholderService interface {
	GetPost(ctx context.Context, id int) (*Post, error)
	GetPosts(ctx context.Context) ([]*Post, error)
	GetComments(ctx context.Context, postID int) ([]*Comment, error)
	GetUser(ctx context.Context, id int) (*User, error)
}

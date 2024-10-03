package jsonplaceholder_impl

import (
	"context"
	"encoding/json"
	"fmt"
	jsonplaceholder "github.com/t34-dev/go-svc-starter/internal/client/json_placeholder"
	"io"
	"net/http"
	"time"
)

// service implements the JSONPlaceholderService interface
type service struct {
	client      *http.Client
	apiKey      string
	userAgent   string
	rateLimiter *time.Ticker
}

// NewService creates a new instance of JSONPlaceholderService
func NewService(options ...ServiceOption) jsonplaceholder.JSONPlaceholderService {
	config := &ServiceConfig{
		Timeout:   time.Second * 30,
		UserAgent: "JSONPlaceholderService/1.0",
	}
	for _, option := range options {
		option(config)
	}

	s := &service{
		client:    newHTTPClient(config),
		apiKey:    config.APIKey,
		userAgent: config.UserAgent,
	}
	if config.RateLimit > 0 {
		s.rateLimiter = time.NewTicker(time.Second / time.Duration(config.RateLimit))
	}
	return s
}

// GetPost retrieves a post by its ID
func (s *service) GetPost(ctx context.Context, id int) (*jsonplaceholder.Post, error) {
	url := fmt.Sprintf("%s/posts/%d", baseURL, id)
	var post jsonplaceholder.Post
	err := s.fetchJSON(ctx, http.MethodGet, url, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts retrieves all posts
func (s *service) GetPosts(ctx context.Context) ([]*jsonplaceholder.Post, error) {
	url := fmt.Sprintf("%s/posts", baseURL)
	var posts []*jsonplaceholder.Post
	err := s.fetchJSON(ctx, http.MethodGet, url, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetComments retrieves comments for a post
func (s *service) GetComments(ctx context.Context, postID int) ([]*jsonplaceholder.Comment, error) {
	url := fmt.Sprintf("%s/posts/%d/comments", baseURL, postID)
	var comments []*jsonplaceholder.Comment
	err := s.fetchJSON(ctx, http.MethodGet, url, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetUser retrieves user information by ID
func (s *service) GetUser(ctx context.Context, id int) (*jsonplaceholder.User, error) {
	url := fmt.Sprintf("%s/users/%d", baseURL, id)
	var user jsonplaceholder.User
	err := s.fetchJSON(ctx, http.MethodGet, url, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// fetchJSON performs an HTTP request and decodes the JSON response
func (s *service) fetchJSON(ctx context.Context, method, url string, target interface{}) error {
	if s.rateLimiter != nil {
		select {
		case <-s.rateLimiter.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

package common

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
)

func (s *ImplementedCommon) GetPost(ctx context.Context, req *common_v1.PostRequest) (*common_v1.PostResponse, error) {
	post, err := s.service.Common.GetPost(ctx, int64(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &common_v1.PostResponse{
		UserId: post.UserId,
		Id:     post.Id,
		Title:  post.Title,
		Body:   post.Body,
	}, nil
}

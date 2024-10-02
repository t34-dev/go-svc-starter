package common_service

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/service"
	"github.com/t34-dev/go-svc-starter/internal/tracing"
	"github.com/t34-dev/go-svc-starter/internal/validator"
	"github.com/t34-dev/go-utils/pkg/sys/validate"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"
)

var _ service.CommonService = &commonService{}

type commonService struct {
	opt service.Dependencies
}

func New(opt service.Dependencies) service.CommonService {
	return &commonService{
		opt: opt,
	}
}

func (s *commonService) GetDBTime(ctx context.Context) (time.Time, error) {
	ctx, finish := tracing.TraceFunc(ctx, "commonService.GetDBTime", nil)
	defer finish()

	return s.opt.Repos.Common.GetTime(ctx)
}

func (s *commonService) GetTime(ctx context.Context) (time.Time, error) {
	ctx, finish := tracing.TraceFunc(ctx, "commonService.GetTime", nil)
	defer finish()
	return time.Now(), nil
}

func (s *commonService) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	ctx, finish := tracing.TraceFunc(ctx, "commonService.GetPost", map[string]interface{}{"id": id})
	defer finish()

	// validate
	err := validate.Validate(ctx, validator.ValidateID(id))
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// parse JSON using gjson
	result := gjson.ParseBytes(body)

	return &model.Post{
		UserId: result.Get("userId").Int(),
		Id:     result.Get("id").Int(),
		Title:  result.Get("title").String(),
		Body:   result.Get("body").String(),
	}, nil
}

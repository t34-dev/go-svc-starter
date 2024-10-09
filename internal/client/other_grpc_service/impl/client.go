package othergrpcservice_impl

import (
	"context"
	othergrpcservice "github.com/t34-dev/go-svc-starter/internal/client/other_grpc_service"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/trace"
)

const isTSL = false

var _ othergrpcservice.OtherGRPCService = (*client)(nil)

type client struct {
	commonClient common_v1.CommonV1Client
}

func New(commonClient common_v1.CommonV1Client) othergrpcservice.OtherGRPCService {
	return &client{commonClient: commonClient}
}

func (c client) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	ctx, finish := trace.TraceFunc(ctx, "client.GetPost", map[string]interface{}{"id": id})
	defer finish()

	res, err := c.commonClient.GetPost(ctx, &common_v1.PostRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &model.Post{
		UserId: res.GetUserId(),
		Id:     res.GetId(),
		Title:  res.GetTitle(),
		Body:   res.GetBody(),
	}, nil
}

func (c client) GetPosts(ctx context.Context) ([]*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

//func New(noteClient desc.OtherNoteV1Client) othergrpcservice.OtherGRPCService {
//	return &client{noteClient: noteClient}
//}

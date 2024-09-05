// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.3
// source: random_v1/random.proto

package random_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RandomServiceClient is the client API for RandomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RandomServiceClient interface {
	GetPing(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PongResponse, error)
	// Get current time
	GetCurrentTime(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*TimeResponse, error)
	// Get a random number
	GetRandomNumber(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*NumberResponse, error)
	// Get a random quote
	GetRandomQuote(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*QuoteResponse, error)
	// Perform a long operation with streaming response
	PerformLongOperation(ctx context.Context, in *LongOperationRequest, opts ...grpc.CallOption) (RandomService_PerformLongOperationClient, error)
	// Get text length
	GetLen(ctx context.Context, in *TxtRequest, opts ...grpc.CallOption) (*TxtResponse, error)
	// Get person information
	GetPerson(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Person, error)
}

type randomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRandomServiceClient(cc grpc.ClientConnInterface) RandomServiceClient {
	return &randomServiceClient{cc}
}

func (c *randomServiceClient) GetPing(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *randomServiceClient) GetCurrentTime(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*TimeResponse, error) {
	out := new(TimeResponse)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetCurrentTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *randomServiceClient) GetRandomNumber(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*NumberResponse, error) {
	out := new(NumberResponse)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetRandomNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *randomServiceClient) GetRandomQuote(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*QuoteResponse, error) {
	out := new(QuoteResponse)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetRandomQuote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *randomServiceClient) PerformLongOperation(ctx context.Context, in *LongOperationRequest, opts ...grpc.CallOption) (RandomService_PerformLongOperationClient, error) {
	stream, err := c.cc.NewStream(ctx, &RandomService_ServiceDesc.Streams[0], "/random_v1.RandomService/PerformLongOperation", opts...)
	if err != nil {
		return nil, err
	}
	x := &randomServicePerformLongOperationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RandomService_PerformLongOperationClient interface {
	Recv() (*LongOperationResponse, error)
	grpc.ClientStream
}

type randomServicePerformLongOperationClient struct {
	grpc.ClientStream
}

func (x *randomServicePerformLongOperationClient) Recv() (*LongOperationResponse, error) {
	m := new(LongOperationResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *randomServiceClient) GetLen(ctx context.Context, in *TxtRequest, opts ...grpc.CallOption) (*TxtResponse, error) {
	out := new(TxtResponse)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetLen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *randomServiceClient) GetPerson(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Person, error) {
	out := new(Person)
	err := c.cc.Invoke(ctx, "/random_v1.RandomService/GetPerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RandomServiceServer is the server API for RandomService service.
// All implementations must embed UnimplementedRandomServiceServer
// for forward compatibility
type RandomServiceServer interface {
	GetPing(context.Context, *EmptyRequest) (*PongResponse, error)
	// Get current time
	GetCurrentTime(context.Context, *EmptyRequest) (*TimeResponse, error)
	// Get a random number
	GetRandomNumber(context.Context, *EmptyRequest) (*NumberResponse, error)
	// Get a random quote
	GetRandomQuote(context.Context, *EmptyRequest) (*QuoteResponse, error)
	// Perform a long operation with streaming response
	PerformLongOperation(*LongOperationRequest, RandomService_PerformLongOperationServer) error
	// Get text length
	GetLen(context.Context, *TxtRequest) (*TxtResponse, error)
	// Get person information
	GetPerson(context.Context, *EmptyRequest) (*Person, error)
	mustEmbedUnimplementedRandomServiceServer()
}

// UnimplementedRandomServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRandomServiceServer struct {
}

func (UnimplementedRandomServiceServer) GetPing(context.Context, *EmptyRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPing not implemented")
}
func (UnimplementedRandomServiceServer) GetCurrentTime(context.Context, *EmptyRequest) (*TimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentTime not implemented")
}
func (UnimplementedRandomServiceServer) GetRandomNumber(context.Context, *EmptyRequest) (*NumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRandomNumber not implemented")
}
func (UnimplementedRandomServiceServer) GetRandomQuote(context.Context, *EmptyRequest) (*QuoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRandomQuote not implemented")
}
func (UnimplementedRandomServiceServer) PerformLongOperation(*LongOperationRequest, RandomService_PerformLongOperationServer) error {
	return status.Errorf(codes.Unimplemented, "method PerformLongOperation not implemented")
}
func (UnimplementedRandomServiceServer) GetLen(context.Context, *TxtRequest) (*TxtResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLen not implemented")
}
func (UnimplementedRandomServiceServer) GetPerson(context.Context, *EmptyRequest) (*Person, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPerson not implemented")
}
func (UnimplementedRandomServiceServer) mustEmbedUnimplementedRandomServiceServer() {}

// UnsafeRandomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RandomServiceServer will
// result in compilation errors.
type UnsafeRandomServiceServer interface {
	mustEmbedUnimplementedRandomServiceServer()
}

func RegisterRandomServiceServer(s grpc.ServiceRegistrar, srv RandomServiceServer) {
	s.RegisterService(&RandomService_ServiceDesc, srv)
}

func _RandomService_GetPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetPing(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RandomService_GetCurrentTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetCurrentTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetCurrentTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetCurrentTime(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RandomService_GetRandomNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetRandomNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetRandomNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetRandomNumber(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RandomService_GetRandomQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetRandomQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetRandomQuote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetRandomQuote(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RandomService_PerformLongOperation_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LongOperationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RandomServiceServer).PerformLongOperation(m, &randomServicePerformLongOperationServer{stream})
}

type RandomService_PerformLongOperationServer interface {
	Send(*LongOperationResponse) error
	grpc.ServerStream
}

type randomServicePerformLongOperationServer struct {
	grpc.ServerStream
}

func (x *randomServicePerformLongOperationServer) Send(m *LongOperationResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _RandomService_GetLen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxtRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetLen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetLen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetLen(ctx, req.(*TxtRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RandomService_GetPerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RandomServiceServer).GetPerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/random_v1.RandomService/GetPerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RandomServiceServer).GetPerson(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RandomService_ServiceDesc is the grpc.ServiceDesc for RandomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RandomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "random_v1.RandomService",
	HandlerType: (*RandomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPing",
			Handler:    _RandomService_GetPing_Handler,
		},
		{
			MethodName: "GetCurrentTime",
			Handler:    _RandomService_GetCurrentTime_Handler,
		},
		{
			MethodName: "GetRandomNumber",
			Handler:    _RandomService_GetRandomNumber_Handler,
		},
		{
			MethodName: "GetRandomQuote",
			Handler:    _RandomService_GetRandomQuote_Handler,
		},
		{
			MethodName: "GetLen",
			Handler:    _RandomService_GetLen_Handler,
		},
		{
			MethodName: "GetPerson",
			Handler:    _RandomService_GetPerson_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PerformLongOperation",
			Handler:       _RandomService_PerformLongOperation_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "random_v1/random.proto",
}

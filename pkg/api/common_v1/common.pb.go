// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.3
// source: common_v1/common.proto

package common_v1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LongOperationResponse_Status int32

const (
	LongOperationResponse_IN_PROGRESS LongOperationResponse_Status = 0
	LongOperationResponse_COMPLETED   LongOperationResponse_Status = 1
	LongOperationResponse_FAILED      LongOperationResponse_Status = 2
)

// Enum value maps for LongOperationResponse_Status.
var (
	LongOperationResponse_Status_name = map[int32]string{
		0: "IN_PROGRESS",
		1: "COMPLETED",
		2: "FAILED",
	}
	LongOperationResponse_Status_value = map[string]int32{
		"IN_PROGRESS": 0,
		"COMPLETED":   1,
		"FAILED":      2,
	}
)

func (x LongOperationResponse_Status) Enum() *LongOperationResponse_Status {
	p := new(LongOperationResponse_Status)
	*p = x
	return p
}

func (x LongOperationResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LongOperationResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_common_proto_enumTypes[0].Descriptor()
}

func (LongOperationResponse_Status) Type() protoreflect.EnumType {
	return &file_common_v1_common_proto_enumTypes[0]
}

func (x LongOperationResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LongOperationResponse_Status.Descriptor instead.
func (LongOperationResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_common_proto_rawDescGZIP(), []int{2, 0}
}

type TimeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *TimeResponse) Reset() {
	*x = TimeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeResponse) ProtoMessage() {}

func (x *TimeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeResponse.ProtoReflect.Descriptor instead.
func (*TimeResponse) Descriptor() ([]byte, []int) {
	return file_common_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *TimeResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type LongOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LongOperationRequest) Reset() {
	*x = LongOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongOperationRequest) ProtoMessage() {}

func (x *LongOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongOperationRequest.ProtoReflect.Descriptor instead.
func (*LongOperationRequest) Descriptor() ([]byte, []int) {
	return file_common_v1_common_proto_rawDescGZIP(), []int{1}
}

type LongOperationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status   LongOperationResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=common_v1.LongOperationResponse_Status" json:"status,omitempty"`
	Message  string                       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Progress int32                        `protobuf:"varint,3,opt,name=progress,proto3" json:"progress,omitempty"`
	Result   string                       `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *LongOperationResponse) Reset() {
	*x = LongOperationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongOperationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongOperationResponse) ProtoMessage() {}

func (x *LongOperationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongOperationResponse.ProtoReflect.Descriptor instead.
func (*LongOperationResponse) Descriptor() ([]byte, []int) {
	return file_common_v1_common_proto_rawDescGZIP(), []int{2}
}

func (x *LongOperationResponse) GetStatus() LongOperationResponse_Status {
	if x != nil {
		return x.Status
	}
	return LongOperationResponse_IN_PROGRESS
}

func (x *LongOperationResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *LongOperationResponse) GetProgress() int32 {
	if x != nil {
		return x.Progress
	}
	return 0
}

func (x *LongOperationResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_common_v1_common_proto protoreflect.FileDescriptor

var file_common_v1_common_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x69, 0x0a, 0x0c, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x59, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x29, 0x92, 0x41, 0x26, 0x32, 0x0c,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x4a, 0x16, 0x22, 0x32,
	0x30, 0x32, 0x33, 0x2d, 0x30, 0x34, 0x2d, 0x30, 0x31, 0x54, 0x31, 0x32, 0x3a, 0x30, 0x30, 0x3a,
	0x30, 0x30, 0x5a, 0x22, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x4c, 0x6f,
	0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0xd4, 0x02, 0x0a, 0x15, 0x4c, 0x6f, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x6e, 0x67, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x15, 0x92, 0x41, 0x12, 0x32, 0x10, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x15, 0x92, 0x41, 0x12, 0x32, 0x10, 0x50, 0x72, 0x6f, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42, 0x22, 0x92, 0x41, 0x1f, 0x32, 0x1d, 0x43, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x61, 0x67, 0x65, 0x20, 0x28, 0x30, 0x2d, 0x31, 0x30, 0x30, 0x29, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3c, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x24, 0x92, 0x41, 0x21, 0x32, 0x1f, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x20, 0x28, 0x69, 0x66,
	0x20, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x29, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x34, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a,
	0x0b, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0d,
	0x0a, 0x09, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x32, 0x86, 0x03, 0x0a, 0x08, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x56, 0x31, 0x12, 0x8f, 0x01, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x53, 0x92, 0x41, 0x33, 0x12, 0x10, 0x47, 0x65, 0x74, 0x20, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x1a, 0x1f, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x20,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x17, 0x12, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x56, 0x31,
	0x2f, 0x47, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0xe7, 0x01, 0x0a, 0x0d, 0x4c, 0x6f, 0x6e,
	0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x90, 0x01,
	0x92, 0x41, 0x67, 0x12, 0x18, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x20, 0x61, 0x20, 0x6c,
	0x6f, 0x6e, 0x67, 0x20, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x4b, 0x49,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x73, 0x20, 0x61, 0x20, 0x6c, 0x6f, 0x6e, 0x67, 0x20,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x73, 0x20, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x20,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x20, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x20, 0x69, 0x74,
	0x73, 0x20, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x20,
	0x22, 0x1b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x56, 0x31, 0x2f,
	0x4c, 0x6f, 0x6e, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x01, 0x2a,
	0x30, 0x01, 0x42, 0xde, 0x01, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x33, 0x34, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x76, 0x63,
	0x2d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5f, 0x76, 0x31, 0x92, 0x41, 0x9b, 0x01, 0x12, 0x3e, 0x0a, 0x14, 0x53, 0x53, 0x4f, 0x20,
	0x54, 0x65, 0x73, 0x74, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x41, 0x50, 0x49,
	0x12, 0x21, 0x41, 0x50, 0x49, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65,
	0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x02, 0x01, 0x02, 0x32, 0x10, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x3a, 0x10,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x5a, 0x1f, 0x0a, 0x1d, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x0f, 0x08, 0x02, 0x1a, 0x09, 0x58, 0x2d, 0x41, 0x50, 0x49, 0x2d, 0x4b, 0x65, 0x79, 0x20,
	0x02, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74,
	0x68, 0x12, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_v1_common_proto_rawDescOnce sync.Once
	file_common_v1_common_proto_rawDescData = file_common_v1_common_proto_rawDesc
)

func file_common_v1_common_proto_rawDescGZIP() []byte {
	file_common_v1_common_proto_rawDescOnce.Do(func() {
		file_common_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_common_proto_rawDescData)
	})
	return file_common_v1_common_proto_rawDescData
}

var file_common_v1_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_common_v1_common_proto_goTypes = []interface{}{
	(LongOperationResponse_Status)(0), // 0: common_v1.LongOperationResponse.Status
	(*TimeResponse)(nil),              // 1: common_v1.TimeResponse
	(*LongOperationRequest)(nil),      // 2: common_v1.LongOperationRequest
	(*LongOperationResponse)(nil),     // 3: common_v1.LongOperationResponse
	(*timestamppb.Timestamp)(nil),     // 4: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),             // 5: google.protobuf.Empty
}
var file_common_v1_common_proto_depIdxs = []int32{
	4, // 0: common_v1.TimeResponse.time:type_name -> google.protobuf.Timestamp
	0, // 1: common_v1.LongOperationResponse.status:type_name -> common_v1.LongOperationResponse.Status
	5, // 2: common_v1.CommonV1.GetTime:input_type -> google.protobuf.Empty
	2, // 3: common_v1.CommonV1.LongOperation:input_type -> common_v1.LongOperationRequest
	1, // 4: common_v1.CommonV1.GetTime:output_type -> common_v1.TimeResponse
	3, // 5: common_v1.CommonV1.LongOperation:output_type -> common_v1.LongOperationResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_v1_common_proto_init() }
func file_common_v1_common_proto_init() {
	if File_common_v1_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_v1_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_v1_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LongOperationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_v1_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LongOperationResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_v1_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_common_v1_common_proto_goTypes,
		DependencyIndexes: file_common_v1_common_proto_depIdxs,
		EnumInfos:         file_common_v1_common_proto_enumTypes,
		MessageInfos:      file_common_v1_common_proto_msgTypes,
	}.Build()
	File_common_v1_common_proto = out.File
	file_common_v1_common_proto_rawDesc = nil
	file_common_v1_common_proto_goTypes = nil
	file_common_v1_common_proto_depIdxs = nil
}

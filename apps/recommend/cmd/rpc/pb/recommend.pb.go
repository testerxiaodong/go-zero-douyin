// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: rpc/pb/recommend.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VideoRecommendSectionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id"`
	SectionId int64 `protobuf:"varint,2,opt,name=section_id,json=sectionId,proto3" json:"section_id"`
	Count     int64 `protobuf:"varint,3,opt,name=count,proto3" json:"count"`
}

func (x *VideoRecommendSectionReq) Reset() {
	*x = VideoRecommendSectionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_pb_recommend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VideoRecommendSectionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoRecommendSectionReq) ProtoMessage() {}

func (x *VideoRecommendSectionReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_pb_recommend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoRecommendSectionReq.ProtoReflect.Descriptor instead.
func (*VideoRecommendSectionReq) Descriptor() ([]byte, []int) {
	return file_rpc_pb_recommend_proto_rawDescGZIP(), []int{0}
}

func (x *VideoRecommendSectionReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *VideoRecommendSectionReq) GetSectionId() int64 {
	if x != nil {
		return x.SectionId
	}
	return 0
}

func (x *VideoRecommendSectionReq) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type VideoRecommendSectionResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoIds []int64 `protobuf:"varint,1,rep,packed,name=video_ids,json=videoIds,proto3" json:"video_ids"`
}

func (x *VideoRecommendSectionResp) Reset() {
	*x = VideoRecommendSectionResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_pb_recommend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VideoRecommendSectionResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoRecommendSectionResp) ProtoMessage() {}

func (x *VideoRecommendSectionResp) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_pb_recommend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoRecommendSectionResp.ProtoReflect.Descriptor instead.
func (*VideoRecommendSectionResp) Descriptor() ([]byte, []int) {
	return file_rpc_pb_recommend_proto_rawDescGZIP(), []int{1}
}

func (x *VideoRecommendSectionResp) GetVideoIds() []int64 {
	if x != nil {
		return x.VideoIds
	}
	return nil
}

var File_rpc_pb_recommend_proto protoreflect.FileDescriptor

var file_rpc_pb_recommend_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x68, 0x0a, 0x18,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x38, 0x0a, 0x19, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52,
	0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x08, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x73,
	0x32, 0x61, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x12, 0x54, 0x0a,
	0x15, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x69, 0x64, 0x65,
	0x6f, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52,
	0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_rpc_pb_recommend_proto_rawDescOnce sync.Once
	file_rpc_pb_recommend_proto_rawDescData = file_rpc_pb_recommend_proto_rawDesc
)

func file_rpc_pb_recommend_proto_rawDescGZIP() []byte {
	file_rpc_pb_recommend_proto_rawDescOnce.Do(func() {
		file_rpc_pb_recommend_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_pb_recommend_proto_rawDescData)
	})
	return file_rpc_pb_recommend_proto_rawDescData
}

var file_rpc_pb_recommend_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_pb_recommend_proto_goTypes = []interface{}{
	(*VideoRecommendSectionReq)(nil),  // 0: pb.VideoRecommendSectionReq
	(*VideoRecommendSectionResp)(nil), // 1: pb.VideoRecommendSectionResp
}
var file_rpc_pb_recommend_proto_depIdxs = []int32{
	0, // 0: pb.recommend.VideoRecommendSection:input_type -> pb.VideoRecommendSectionReq
	1, // 1: pb.recommend.VideoRecommendSection:output_type -> pb.VideoRecommendSectionResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_pb_recommend_proto_init() }
func file_rpc_pb_recommend_proto_init() {
	if File_rpc_pb_recommend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_pb_recommend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VideoRecommendSectionReq); i {
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
		file_rpc_pb_recommend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VideoRecommendSectionResp); i {
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
			RawDescriptor: file_rpc_pb_recommend_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_pb_recommend_proto_goTypes,
		DependencyIndexes: file_rpc_pb_recommend_proto_depIdxs,
		MessageInfos:      file_rpc_pb_recommend_proto_msgTypes,
	}.Build()
	File_rpc_pb_recommend_proto = out.File
	file_rpc_pb_recommend_proto_rawDesc = nil
	file_rpc_pb_recommend_proto_goTypes = nil
	file_rpc_pb_recommend_proto_depIdxs = nil
}

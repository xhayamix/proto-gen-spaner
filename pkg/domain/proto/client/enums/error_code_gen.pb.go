// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.26.1
// source: client/enums/error_code_gen.proto

package enums

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

type ErrorCode int32

const (
	ErrorCode_ErrorCode_Unknown ErrorCode = 0
	// パラメータの不正
	ErrorCode_ErrorCode_InvalidArgument ErrorCode = 1001
	// サーバー内部エラー
	ErrorCode_ErrorCode_Internal ErrorCode = 1002
	// 認証エラー
	ErrorCode_ErrorCode_Unauthenticated ErrorCode = 1003
	// アクセス権限なし
	ErrorCode_ErrorCode_PermissionDenied ErrorCode = 1004
	// リソースが見つからなかった
	ErrorCode_ErrorCode_NotFound ErrorCode = 1005
	// ユーザーが見つからなかった
	ErrorCode_ErrorCode_UserNotFound ErrorCode = 2001
	// ユーザーは削除済み
	ErrorCode_ErrorCode_UserDeleted ErrorCode = 2002
	// メンテナンス中
	ErrorCode_ErrorCode_InMaintenance ErrorCode = 2003
	// アカウント停止中
	ErrorCode_ErrorCode_AccountBan ErrorCode = 2004
	// NGワードが含まれている
	ErrorCode_ErrorCode_NgWordContains ErrorCode = 2005
	// 無効な日付
	ErrorCode_ErrorCode_ShopInvalidDay ErrorCode = 2006
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0:    "ErrorCode_Unknown",
		1001: "ErrorCode_InvalidArgument",
		1002: "ErrorCode_Internal",
		1003: "ErrorCode_Unauthenticated",
		1004: "ErrorCode_PermissionDenied",
		1005: "ErrorCode_NotFound",
		2001: "ErrorCode_UserNotFound",
		2002: "ErrorCode_UserDeleted",
		2003: "ErrorCode_InMaintenance",
		2004: "ErrorCode_AccountBan",
		2005: "ErrorCode_NgWordContains",
		2006: "ErrorCode_ShopInvalidDay",
	}
	ErrorCode_value = map[string]int32{
		"ErrorCode_Unknown":          0,
		"ErrorCode_InvalidArgument":  1001,
		"ErrorCode_Internal":         1002,
		"ErrorCode_Unauthenticated":  1003,
		"ErrorCode_PermissionDenied": 1004,
		"ErrorCode_NotFound":         1005,
		"ErrorCode_UserNotFound":     2001,
		"ErrorCode_UserDeleted":      2002,
		"ErrorCode_InMaintenance":    2003,
		"ErrorCode_AccountBan":       2004,
		"ErrorCode_NgWordContains":   2005,
		"ErrorCode_ShopInvalidDay":   2006,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_client_enums_error_code_gen_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_client_enums_error_code_gen_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_client_enums_error_code_gen_proto_rawDescGZIP(), []int{0}
}

var File_client_enums_error_code_gen_proto protoreflect.FileDescriptor

var file_client_enums_error_code_gen_proto_rawDesc = []byte{
	0x0a, 0x21, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x65, 0x6e, 0x75, 0x6d,
	0x73, 0x2a, 0xe5, 0x02, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x15, 0x0a, 0x11, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x55, 0x6e, 0x6b,
	0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x19, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x5f, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x41, 0x72, 0x67, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x10, 0xe9, 0x07, 0x12, 0x17, 0x0a, 0x12, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x5f, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x10, 0xea, 0x07, 0x12,
	0x1e, 0x0a, 0x19, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x55, 0x6e, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x10, 0xeb, 0x07, 0x12,
	0x1f, 0x0a, 0x1a, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0xec, 0x07,
	0x12, 0x17, 0x0a, 0x12, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x4e, 0x6f,
	0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xed, 0x07, 0x12, 0x1b, 0x0a, 0x16, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x46, 0x6f,
	0x75, 0x6e, 0x64, 0x10, 0xd1, 0x0f, 0x12, 0x1a, 0x0a, 0x15, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x10,
	0xd2, 0x0f, 0x12, 0x1c, 0x0a, 0x17, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f,
	0x49, 0x6e, 0x4d, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x10, 0xd3, 0x0f,
	0x12, 0x19, 0x0a, 0x14, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6e, 0x10, 0xd4, 0x0f, 0x12, 0x1d, 0x0a, 0x18, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x4e, 0x67, 0x57, 0x6f, 0x72, 0x64, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x10, 0xd5, 0x0f, 0x12, 0x1d, 0x0a, 0x18, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x5f, 0x53, 0x68, 0x6f, 0x70, 0x49, 0x6e, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x44, 0x61, 0x79, 0x10, 0xd6, 0x0f, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x68, 0x61, 0x79, 0x61, 0x6d, 0x69, 0x78,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x73, 0x70, 0x61, 0x6e, 0x6e,
	0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_enums_error_code_gen_proto_rawDescOnce sync.Once
	file_client_enums_error_code_gen_proto_rawDescData = file_client_enums_error_code_gen_proto_rawDesc
)

func file_client_enums_error_code_gen_proto_rawDescGZIP() []byte {
	file_client_enums_error_code_gen_proto_rawDescOnce.Do(func() {
		file_client_enums_error_code_gen_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_enums_error_code_gen_proto_rawDescData)
	})
	return file_client_enums_error_code_gen_proto_rawDescData
}

var file_client_enums_error_code_gen_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_client_enums_error_code_gen_proto_goTypes = []any{
	(ErrorCode)(0), // 0: client.enums.ErrorCode
}
var file_client_enums_error_code_gen_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_client_enums_error_code_gen_proto_init() }
func file_client_enums_error_code_gen_proto_init() {
	if File_client_enums_error_code_gen_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_client_enums_error_code_gen_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_enums_error_code_gen_proto_goTypes,
		DependencyIndexes: file_client_enums_error_code_gen_proto_depIdxs,
		EnumInfos:         file_client_enums_error_code_gen_proto_enumTypes,
	}.Build()
	File_client_enums_error_code_gen_proto = out.File
	file_client_enums_error_code_gen_proto_rawDesc = nil
	file_client_enums_error_code_gen_proto_goTypes = nil
	file_client_enums_error_code_gen_proto_depIdxs = nil
}

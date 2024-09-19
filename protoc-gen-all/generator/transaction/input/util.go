package input

import (
	"fmt"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

var ServerMessageAccessorSet = NewMessageAccessorSet(
	MessageAccessorType_AdminAndServer,
	MessageAccessorType_All,
	MessageAccessorType_AllWithCommonResponse,
)

var ClientMessageAccessorSet = NewMessageAccessorSet(
	MessageAccessorType_OnlyClient,
	MessageAccessorType_OnlyClientWithCommonResponse,
	MessageAccessorType_All,
	MessageAccessorType_AllWithCommonResponse,
)

var ClientMessageCommonResponseAccessorSet = NewMessageAccessorSet(
	MessageAccessorType_OnlyClientWithCommonResponse,
	MessageAccessorType_AllWithCommonResponse,
)

var AdminFieldAccessorSet = NewFieldAccessorSet(
	FieldAccessorType_All,
	FieldAccessorType_OnlyAdmin,
	FieldAccessorType_AdminAndServer,
	FieldAccessorType_AdminAndClient,
)

var ServerFieldAccessorSet = NewFieldAccessorSet(
	FieldAccessorType_All,
	FieldAccessorType_OnlyServer,
	FieldAccessorType_AdminAndServer,
	FieldAccessorType_ServerAndClient,
)

var ClientFieldAccessorSet = NewFieldAccessorSet(
	FieldAccessorType_All,
	FieldAccessorType_OnlyClient,
	FieldAccessorType_AdminAndClient,
	FieldAccessorType_ServerAndClient,
)

type spannerType = string

const (
	spannerTypeBool   spannerType = "BOOL"
	spannerTypeInt    spannerType = "INT64"
	spannerTypeString spannerType = "STRING(MAX)"
	spannerTypeBytes  spannerType = "BYTES(MAX)"
	spannerTypeTime   spannerType = "TIMESTAMP"
)

func GetSpannerType(typeKind TypeKind, fieldSnakeName string, isList bool) (spannerType, error) {
	var typeName string
	switch typeKind {
	case TypeKind_Bool:
		typeName = spannerTypeBool
	case TypeKind_Int32:
		typeName = spannerTypeInt
	case TypeKind_Int64:
		typeName = spannerTypeInt
	case TypeKind_String:
		typeName = spannerTypeString
	case TypeKind_Enum:
		typeName = spannerTypeInt
	case TypeKind_Bytes:
		typeName = spannerTypeBytes
	default:
		return "", perrors.Newf("サポートされていないTypeKindです。 TypeKind = %v", typeKind)
	}
	if core.IsTimeField(fieldSnakeName) {
		typeName = spannerTypeTime
	}
	if isList {
		typeName = fmt.Sprintf("ARRAY<%s>", typeName)
	}
	return typeName, nil
}

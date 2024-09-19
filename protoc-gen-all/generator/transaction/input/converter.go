package input

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	options "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/server/options/transaction"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

func ConvertMessageFromProto(file *protogen.File) (*Message, error) {
	if len(file.Messages) != 1 {
		return nil, perrors.Newf("このprotoファイルはメッセージの定義数を1にする必要があります。 file = %v", file.Desc.FullName())
	}

	return convert(file.Messages[0])
}

func convert(message *protogen.Message) (*Message, error) {
	messageOption, ok := proto.GetExtension(message.Desc.Options(), options.E_Message).(*options.MessageOption)
	if !ok {
		return nil, perrors.Newf("型アサーションに失敗しました")
	}

	messageAccessorType, err := ConvertMessageAccessorTypeFromProto(messageOption.GetAccessorType())
	if err != nil {
		return nil, perrors.Stack(err)
	}

	var indexes []*Index
	var interleave *Interleave
	var ttl *TTL
	if messageOption.GetDdl() != nil {
		indexes = make([]*Index, 0, len(messageOption.GetDdl().GetIndexes()))
		for _, index := range messageOption.GetDdl().GetIndexes() {
			keys := make([]*IndexKey, 0, len(index.GetKeys()))
			for _, key := range index.GetKeys() {
				keys = append(keys, &IndexKey{
					SnakeName: core.ToSnakeCase(key.Column),
					Desc:      key.Desc,
				})
			}
			storing := make([]string, 0, len(index.GetStoring()))
			for _, s := range index.GetStoring() {
				storing = append(storing, core.ToSnakeCase(s))
			}

			indexes = append(indexes, &Index{
				Keys:         keys,
				Unique:       index.Unique,
				SnakeStoring: storing,
			})
		}

		if messageOption.GetDdl().GetInterleave() != nil {
			interleave = &Interleave{
				TableSnakeName: core.ToSnakeCase(messageOption.GetDdl().GetInterleave().GetTable()),
			}
		}

		if messageOption.GetDdl().GetTtl() != nil {
			protoTTL := messageOption.GetDdl().GetTtl()
			ttl = &TTL{
				TimestampColumnSnakeName: protoTTL.GetTimestampColumn(),
				Days:                     protoTTL.GetDays(),
			}
			if ttl.TimestampColumnSnakeName == "" {
				ttl.TimestampColumnSnakeName = "updated_time"
			}
		}
	}

	ret := &Message{
		Messages:  make([]*Message, 0, message.Desc.Messages().Len()),
		SnakeName: core.ToSnakeCase(string(message.Desc.FullName().Name())),
		Comment:   core.CommentReplacer.Replace(message.Comments.Leading.String()),
		Fields:    nil,
		Option: &MessageOption{
			AccessorType: messageAccessorType,
			DDL: &MessageOptionDDL{
				Indexes:    indexes,
				Interleave: interleave,
				TTL:        ttl,
			},
			InsertTiming: messageOption.GetInsertTiming(),
		},
	}

	inputFields := make([]*Field, 0, len(message.Fields)+2)
	for _, field := range message.Fields {
		var pkgType PkgType
		var importFilePath string
		var typeName string
		var typeKind TypeKind
		var rawTypeName string
		var rawTypeKind TypeKind

		fieldOption, ok := proto.GetExtension(field.Desc.Options(), options.E_Field).(*options.FieldOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}
		fieldAccessorType, err := ConvertFieldAccessorTypeFromProto(fieldOption.GetAccessorType())
		if err != nil {
			return nil, perrors.Stack(err)
		}

		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			typeName = FieldType_Bool
			typeKind = TypeKind_Bool
			rawTypeName = FieldType_Bool
			rawTypeKind = TypeKind_Bool
		case protoreflect.Int32Kind:
			// Spannerに合わせる
			typeName = FieldType_Int64
			typeKind = TypeKind_Int64
			rawTypeName = FieldType_Int32
			rawTypeKind = TypeKind_Int32
		case protoreflect.Int64Kind:
			typeName = FieldType_Int64
			typeKind = TypeKind_Int64
			rawTypeName = FieldType_Int64
			rawTypeKind = TypeKind_Int64
		case protoreflect.StringKind:
			typeName = FieldType_String
			typeKind = TypeKind_String
			rawTypeName = FieldType_String
			rawTypeKind = TypeKind_String
		case protoreflect.EnumKind:
			typeName = string(field.Desc.Enum().Name())
			typeKind = TypeKind_Enum
			rawTypeName = string(field.Desc.Enum().Name())
			rawTypeKind = TypeKind_Enum
		case protoreflect.BytesKind:
			typeName = FieldType_Bytes
			typeKind = TypeKind_Bytes
			rawTypeName = FieldType_Bytes
			rawTypeKind = TypeKind_Bytes
		case protoreflect.MessageKind:
			// 型がmessageのfieldはAdmin系では使えない
			if AdminFieldAccessorSet.Contains(fieldAccessorType) {
				return nil, perrors.Newf("transactionデータでDB定義がある場合は型がmessageのフィールドを利用できません。 tableName = %s, columnName = %s", ret.SnakeName, core.ToSnakeCase(field.Desc.TextName()))
			}

			pkgName := string(field.Desc.Message().ParentFile().Package())
			if strings.HasSuffix(pkgName, "common") {
				pkgType = PkgType_Common
				importFilePath = field.Desc.Message().ParentFile().Path()
			} else {
				m, err := convert(field.Message)
				if err != nil {
					return nil, perrors.Stack(err)
				}
				ret.Messages = append(ret.Messages, m)
			}
			rawTypeName = string(field.Desc.Message().Name())
			rawTypeKind = TypeKind_Message
		case protoreflect.DoubleKind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.GroupKind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			return nil, perrors.Newf("サポートされていないKindです。 Kind = %v", field.Desc.Kind().String())
		default:
			return nil, perrors.Newf("サポートされていないKindです。 Kind = %v", field.Desc.Kind().String())
		}

		paths := strings.Split(importFilePath, "/")
		fileName := paths[len(paths)-1]
		if importFilePath == field.Desc.ParentFile().Path() {
			fileName = ""
		}

		inputField := &Field{
			PkgType:        pkgType,
			ImportFileName: fileName,
			SnakeName:      core.ToSnakeCase(field.Desc.TextName()),
			Comment:        core.CommentReplacer.Replace(field.Comments.Leading.String()),
			Type:           typeName,
			TypeKind:       typeKind,
			RawType:        rawTypeName,
			RawTypeKind:    rawTypeKind,
			IsList:         field.Desc.IsList(),
			Number:         int32(field.Desc.Number()),
			Option:         nil,
		}

		var pk bool
		var masterRef *MasterRef
		if fieldOption.GetDdl() != nil {
			pk = fieldOption.GetDdl().GetPk()
			if fieldOption.GetDdl().GetMasterRef() != nil {
				fields := fieldOption.GetDdl().GetMasterRef().GetParentColumns()
				parentColumnSnakeNames := make([]string, 0, len(fields))
				for _, f := range fields {
					parentColumnSnakeNames = append(parentColumnSnakeNames, core.ToSnakeCase(f))
				}

				masterRef = &MasterRef{
					TableSnakeName:         core.ToSnakeCase(fieldOption.GetDdl().GetMasterRef().GetTable()),
					ColumnSnakeName:        core.ToSnakeCase(fieldOption.GetDdl().GetMasterRef().GetColumn()),
					ParentColumnSnakeNames: parentColumnSnakeNames,
				}
			}
		}

		option := &FieldOption{
			AccessorType: fieldAccessorType,
			DDL: &FieldOptionDDL{
				PK:        pk,
				MasterRef: masterRef,
			},
		}

		inputField.Option = option

		inputFields = append(inputFields, inputField)
	}
	inputFields = append(inputFields,
		&Field{
			SnakeName: "created_time",
			Comment:   "作成日時",
			Type:      "int64",
			TypeKind:  TypeKind_Int64,
			IsList:    false,
			Option: &FieldOption{
				AccessorType: FieldAccessorType_AdminAndServer,
				DDL: &FieldOptionDDL{
					PK:        false,
					MasterRef: nil,
				},
			},
		},
		&Field{
			SnakeName: "updated_time",
			Comment:   "更新日時",
			Type:      "int64",
			TypeKind:  TypeKind_Int64,
			IsList:    false,
			Option: &FieldOption{
				AccessorType: FieldAccessorType_AdminAndServer,
				DDL: &FieldOptionDDL{
					PK:        false,
					MasterRef: nil,
				},
			},
		},
	)
	ret.Fields = inputFields

	return ret, nil
}

func ConvertMessageAccessorTypeFromProto(in options.MessageOption_AccessorType) (MessageAccessorType, error) {
	var out MessageAccessorType

	switch in {
	case options.MessageOption_AdminAndServer:
		out = MessageAccessorType_AdminAndServer
	case options.MessageOption_OnlyClient:
		out = MessageAccessorType_OnlyClient
	case options.MessageOption_OnlyClientWithCommonResponse:
		out = MessageAccessorType_OnlyClientWithCommonResponse
	case options.MessageOption_All:
		out = MessageAccessorType_All
	case options.MessageOption_AllWithCommonResponse:
		out = MessageAccessorType_AllWithCommonResponse
	default:
		return 0, perrors.Newf("サポートされていないAccessorTypeです。 AccessorType = %v", in)
	}

	return out, nil
}

func ConvertFieldAccessorTypeFromProto(in options.FieldOption_AccessorType) (FieldAccessorType, error) {
	var out FieldAccessorType

	switch in {
	case options.FieldOption_All:
		out = FieldAccessorType_All
	case options.FieldOption_OnlyAdmin:
		out = FieldAccessorType_OnlyAdmin
	case options.FieldOption_OnlyServer:
		out = FieldAccessorType_OnlyServer
	case options.FieldOption_OnlyClient:
		out = FieldAccessorType_OnlyClient
	case options.FieldOption_AdminAndServer:
		out = FieldAccessorType_AdminAndServer
	case options.FieldOption_AdminAndClient:
		out = FieldAccessorType_AdminAndClient
	case options.FieldOption_ServerAndClient:
		out = FieldAccessorType_ServerAndClient
	default:
		return 0, perrors.Newf("サポートされていないAccessorTypeです。 AccessorType = %v", in)
	}

	return out, nil
}

type EnumMap map[string]*Enum

func (e EnumMap) Merge(other EnumMap) EnumMap {
	if e == nil {
		return other
	}
	if other == nil {
		return e
	}
	for k, v := range other {
		e[k] = v
	}
	return e
}

func ConvertEnumsFromProto(file *protogen.File) map[string]*Enum {
	res := make(map[string]*Enum, len(file.Enums))
	for _, e := range file.Enums {
		members := make([]*Member, 0, len(e.Values))
		commentMap := make(map[int32]string, len(e.Values))
		for _, v := range e.Values {
			number := int32(v.Desc.Number())
			if number == 0 {
				continue
			}
			comment := core.CommentReplacer.Replace(v.Comments.Leading.String())
			members = append(members, &Member{
				Name:    string(v.Desc.Name()),
				Value:   number,
				Comment: comment,
			})
			commentMap[number] = comment
		}
		name := string(e.Desc.Name())
		res[name] = &Enum{
			Name:              name,
			Members:           members,
			CommentMapByValue: commentMap,
		}
	}
	return res
}

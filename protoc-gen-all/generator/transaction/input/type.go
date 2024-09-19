package input

import (
	"github.com/scylladb/go-set/i32set"
)

type MessageAccessorType int32

const (
	MessageAccessorType_AdminAndServer MessageAccessorType = iota
	MessageAccessorType_OnlyClient
	MessageAccessorType_OnlyClientWithCommonResponse
	MessageAccessorType_All
	MessageAccessorType_AllWithCommonResponse
)

type MessageAccessorSet struct {
	*i32set.Set
}

func (a MessageAccessorSet) Contains(accessorType MessageAccessorType) bool {
	return a.Has(int32(accessorType))
}

func NewMessageAccessorSet(accessorTypes ...MessageAccessorType) *MessageAccessorSet {
	s := i32set.NewWithSize(len(accessorTypes))
	for _, accessorType := range accessorTypes {
		s.Add(int32(accessorType))
	}

	return &MessageAccessorSet{s}
}

type FieldAccessorType int32

const (
	// FieldAccessorType_All デフォルト値がALL
	FieldAccessorType_All FieldAccessorType = iota
	FieldAccessorType_OnlyAdmin
	FieldAccessorType_OnlyServer
	FieldAccessorType_OnlyClient
	FieldAccessorType_AdminAndServer
	FieldAccessorType_AdminAndClient
	FieldAccessorType_ServerAndClient
)

type FieldAccessorSet struct {
	*i32set.Set
}

func (a FieldAccessorSet) Contains(accessorType FieldAccessorType) bool {
	return a.Has(int32(accessorType))
}

func NewFieldAccessorSet(accessorTypes ...FieldAccessorType) *FieldAccessorSet {
	s := i32set.NewWithSize(len(accessorTypes))
	for _, accessorType := range accessorTypes {
		s.Add(int32(accessorType))
	}

	return &FieldAccessorSet{s}
}

type PkgType int32

const (
	PkgType_Common PkgType = iota + 1
)

type TypeKind int32

const (
	TypeKind_Bool TypeKind = iota + 1
	TypeKind_Int32
	TypeKind_Int64
	TypeKind_String
	TypeKind_Enum
	TypeKind_Bytes
	// TypeKind_Message クライアント向けのレスポンスにしか使われない
	TypeKind_Message
)

type FieldType = string

const (
	FieldType_Bool   = "bool"
	FieldType_Int32  = "int32"
	FieldType_Int64  = "int64"
	FieldType_String = "string"
	FieldType_Bytes  = "[]byte"
)

type Interleave struct {
	TableSnakeName string
}

type TTL struct {
	TimestampColumnSnakeName string
	Days                     int32
}

type MasterRef struct {
	TableSnakeName         string
	ColumnSnakeName        string
	ParentColumnSnakeNames []string
}

type FieldOptionDDL struct {
	PK bool
	// nilチェックが必要
	MasterRef *MasterRef
}

type FieldOption struct {
	AccessorType FieldAccessorType
	DDL          *FieldOptionDDL
}

type Field struct {
	PkgType PkgType
	// ImportFileName 外部ファイルのフィールドならそのファイル名(hoge.proto)が入り、
	// このフィールドが定義されているファイルと同じファイルのフィールドなら空文字が入る
	ImportFileName string
	SnakeName      string
	Comment        string
	// TypeKind_Enumの場合はEnum名が入る
	Type     FieldType
	TypeKind TypeKind
	// Protoで定義されているType
	RawType FieldType
	// Protoで定義されているTypeKind
	RawTypeKind TypeKind
	IsList      bool
	Number      int32
	Option      *FieldOption
}

type IndexKey struct {
	SnakeName string
	Desc      bool
}

type Index struct {
	Keys         []*IndexKey
	Unique       bool
	SnakeStoring []string
}

type MessageOptionDDL struct {
	Indexes []*Index
	// nilチェックが必要
	Interleave *Interleave
	TTL        *TTL
}

type MessageOption struct {
	AccessorType MessageAccessorType
	DDL          *MessageOptionDDL
	InsertTiming string
}

type Message struct {
	Messages  []*Message
	SnakeName string
	Comment   string
	Fields    []*Field
	Option    *MessageOption
}

type Enum struct {
	Name    string
	Members []*Member

	CommentMapByValue map[int32]string
}

type Member struct {
	Name    string
	Value   int32
	Comment string
}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-spanner" mock_$GOPACKAGE/mock_$GOFILE
package master

type Record interface {
	PK() string
	ToKeyValue() map[string]interface{}
	GetTypeMap() map[string]string
}

type RecordWithOriginEnum interface{}

type Slice interface {
	EachRecord(func(Record) bool)
	ToMapByMasterVersion(masterVersions []int32) map[int32]Slice
	Len() int
}

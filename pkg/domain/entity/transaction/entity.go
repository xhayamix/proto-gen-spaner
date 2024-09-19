package transaction

// Entity entityのインターフェース
type Entity interface {
	GetVals() []interface{}
}

type Entities interface {
	Len() int
	EachRecord(iterator func(Entity) bool)
}

// PK PrimaryKeyのインターフェース
type PK interface {
	Generate() []interface{}
}

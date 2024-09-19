package base

const ParamBaseKey = "p"

type OrderType string

const (
	OrderTypeASC  OrderType = "ASC"
	OrderTypeDESC OrderType = "DESC"
)

type ConditionOperator string

const (
	ConditionOperatorEq ConditionOperator = "="
	ConditionOperatorIn ConditionOperator = "IN"
)

type SearchResultType int32

const (
	SearchResultTypeNotSearched SearchResultType = 0
	SearchResultTypeFound       SearchResultType = 1
	SearchResultTypeNotFound    SearchResultType = 2
)

type OperationType int32

const (
	OperationTypeUnknown OperationType = 0
	OperationTypeInsert  OperationType = 1
	OperationTypeUpdate  OperationType = 2
	OperationTypeDelete  OperationType = 3
)

type OrderPair struct {
	Column    string
	OrderType OrderType
}

type OrderPairs []*OrderPair

package registry

import (
	"reflect"
	// "github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction/repository"
)

var transactionRepositoryExtSet = []interface{}{}

func TransactionRepositoryExtSet() []interface{} {
	result := make([]interface{}, 0, len(transactionRepositorySet))
	extNameMap := make(map[string]interface{}, len(transactionRepositoryExtSet))

	for _, r := range transactionRepositoryExtSet {
		extNameMap[reflect.TypeOf(r).Out(0).Name()] = r
	}

	for _, r := range transactionRepositorySet {
		if ext, ok := extNameMap[reflect.TypeOf(r).Out(0).Name()]; ok {
			result = append(result, ext)
			continue
		}

		result = append(result, r)
	}

	return result
}

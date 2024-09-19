{{ template "autogen_comment" }}
package registry

import (
	"github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction/repository"
)

var transactionRepositorySet = []interface{}{
	{{ range . -}}
		repository.New{{ .GoName }}Repository,
	{{ end -}}
}

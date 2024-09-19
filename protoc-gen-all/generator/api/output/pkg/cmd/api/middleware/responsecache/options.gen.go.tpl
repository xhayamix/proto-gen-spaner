{{ template "autogen_comment" }}
package responsecache

import (
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

var enableResponseCacheMethodMap = map[string]interface{}{
	{{ range .EnableResponseCacheMethods -}}
	"{{ .Name }}": (*api.{{ .Type }})(nil),
	{{ end -}}
}

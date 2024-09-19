{{ template "autogen_comment" }}

package converter

import (
    "github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
    "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/enums"
)

func toInt32Slice(s []int64) []int32 {
    ret := make([]int32, 0, len(s))

    for _, v := range s {
        ret = append(ret, int32(v))
    }

    return ret
}

{{ range .Enums }}
func toProto{{ . }}Slice(s enum.{{ . }}Slice) []enums.{{ . }} {
    ret := make([]enums.{{ . }}, 0, len(s))

    s.Each(func(e enum.Enum) bool {
        ret = append(ret, enums.{{ . }}(e.Int32()))
        return true
    })

    return ret
}
{{ end }}

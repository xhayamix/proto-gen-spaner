{{ template "autogen_comment" }}
package masterconverter

import (
	"strconv"
	"time"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/master"
	"github.com/xhayamix/proto-gen-spanner/pkg/util/strings"
	timeutil "github.com/xhayamix/proto-gen-spanner/pkg/util/time"
)

func set{{ .PascalName }}(msg *master.{{ .PascalName }}, typ enum.{{ .PascalName }}Type, value string) error {
	toBool := func (v string) (bool, error) {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, cerrors.Wrap(err, cerrors.Internal)
		}
		return b, nil
	}
	_ = toBool

	toInt32 := func(v string) (int32, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, cerrors.Wrap(err, cerrors.Internal)
		}
		return int32(i), nil
	}
	_ = toInt32

	toInt64 := func(v string) (int64, error) {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, cerrors.Wrap(err, cerrors.Internal)
		}
		return int64(i), nil
	}
	_ = toInt64

	toInt32Slice := func(v string) ([]int32, error) {
		s, err := strings.SplitCommaToInt32(v)
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		return s, nil
	}
	_ = toInt32Slice

	toInt64Slice := func(v string) ([]int64, error) {
		s, err := strings.SplitCommaToInt64(v)
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		return s, nil
	}
	_ = toInt64Slice

	toString := func(v string) (string, error) {
		return v, nil
	}
	_ = toString

	toStringSlice := func(v string) ([]string, error) {
		return strings.SplitComma(v), nil
	}
	_ = toStringSlice

	toTime := func(v string) (int64, error) {
		t, err := time.ParseInLocation("2006/01/02 15:04:05", v, time.Local)
		if err != nil {
			return 0, cerrors.Wrapf(err, cerrors.Internal, "Enum値が日時形式ではありません。 enumName = \"{{ .PascalName }}\", value = %q", v)
		}
		return timeutil.ToUnixMilli(&t), nil
	}
	_ = toTime

	toNotDefined := func(v string) (interface{}, error) {
		return nil, cerrors.Newf(cerrors.Internal, "{{ .PascalName }}のEnum名が不正です。 v = %q", v)
	}
	_ = toNotDefined

	switch typ {
	{{ $name := .PascalName -}}
	{{ range .Elements -}}
	case enum.{{ $name }}Type_{{ .PascalName }}:
	{{- if .HasClient }}
		v, err := to{{ .PascalSettingType }}(value)
		if err != nil {
			return cerrors.Stack(err)
		}
		msg.{{ .PascalName }} = v
	{{- end }}
	{{ end -}}
	default:
		return cerrors.Newf(cerrors.Internal, "{{ .PascalName }}のEnum名が不正です。 typ = %q", typ)
	}
	return nil
}

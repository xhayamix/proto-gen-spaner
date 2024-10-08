{{ template "autogen_comment" }}
package mcache

import (
	"strconv"
	"time"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
	"github.com/xhayamix/proto-gen-spanner/pkg/util/strings"
)

func (s *{{ .PascalName }}) Set(typ enum.{{ .PascalName }}Type, value string) error {
	switch typ {
	{{ $name := .PascalName -}}
	{{ range .Elements -}}
	case enum.{{ $name }}Type_{{ .PascalName }}:
	{{- if .HasServer }}
		v, err := s.to{{ .PascalSettingType }}(value)
		if err != nil {
			return cerrors.Stack(err)
		}
		s.{{ .PascalName }} = v
	{{- end }}
	{{ end -}}
	default:
		return cerrors.Newf(cerrors.Internal, "{{ .PascalName }}のEnum名が不正です。 typ = %q", typ)
	}
	return nil
}

func (s *{{ .PascalName }}) toBool(v string) (bool, error) {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, cerrors.Wrap(err, cerrors.Internal)
	}
	return b, nil
}

func (s *{{ .PascalName }}) toInt32(v string) (int32, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, cerrors.Wrap(err, cerrors.Internal)
	}
	return int32(i), nil
}

func (s *{{ .PascalName }}) toInt32Slice(v string) ([]int32, error) {
	slice, err := strings.SplitCommaToInt32(v)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	return slice, nil
}

func (s *{{ .PascalName }}) toString(v string) (string, error) {
	return v, nil
}

func (s *{{ .PascalName }}) toStringSlice(v string) ([]string, error) {
	return strings.SplitComma(v), nil
}

func (s *{{ .PascalName }}) toTime(v string) (time.Time, error) {
	t, err := time.ParseInLocation("2006/01/02 15:04:05", v, time.Local)
	if err != nil {
		return time.Time{}, cerrors.Wrapf(err, cerrors.Internal, "Enum値が日時形式ではありません。 enumName = \"{{ .PascalName }}\", value = %q", v)
	}
	return t, nil
}

func (s *{{ .PascalName }}) toNotDefined(v string) (interface{}, error) {
	return nil, cerrors.Newf(cerrors.Internal, "{{ .PascalName }}のEnum名が不正です。 v = %q", v)
}

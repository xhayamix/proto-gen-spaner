{{ template "autogen_comment" }}
package transaction
{{- $goName := .GoName }}

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/constant"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/dto/column"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
)

const (
	{{ .GoName }}TableName    = "{{ .GoName }}"
	{{ .GoName }}Comment      = "{{ .Comment }}"
	{{ .GoName }}InsertTiming = "{{ .InsertTiming }}"
)

// {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .Type }} `json:"{{ .GoName }},omitempty"`
{{- end }}
}

func (e *{{ .GoName }}) GetPK() *{{ .GoName }}PK {
	return &{{ .GoName }}PK{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}) GetVals() []interface{} {
	return []interface{}{
	{{- range .Columns }}
	{{- if not .IsOnlyServer }}
		e.{{ .GoName }},
	{{- end }}
	{{- end }}
	}
}

func (e *{{ .GoName }}) ToKeyValue() map[string]interface{} {
	return map[string]interface{}{
	{{- range .Columns }}
	{{- if not .IsOnlyServer }}
		"{{ .GoName }}": e.{{ .GoName }},
	{{- end }}
	{{- end }}
	}
}

func (e *{{ .GoName }}) GetTypeMap() map[string]string {
	return map[string]string{
{{- range .Columns }}
	{{- if not .IsOnlyServer }}
		"{{ .GoName }}": "{{ .Type }}",
	{{- end }}
{{- end }}
	}
}

func (e *{{ .GoName }}) SetKeyValue(columns []string, entity []interface{}) []string {
	errs := make([]string, 0, len(columns))
	for index, column := range columns {
		if len(entity) <= index {
			break
		}
		value := entity[index]
		switch column {
		{{- range .Columns }}
		case "{{ .GoName }}":
			{{- if eq .Type "int64" }}
			switch v := value.(type) {
			case int64:
				e.{{ .GoName }} = v
			case int32:
				e.{{ .GoName }} = int64(v)
			case float32:
				e.{{ .GoName }} = int64(v)
			case float64:
				e.{{ .GoName }} = int64(v)
			default:
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if eq .Type "[]int64" }}
			switch v := value.(type) {
			case []int64:
				e.{{ .GoName }} = v
			case []int32:
				e.{{ .GoName }} = make([]int64, 0, len(v))
				for _, i := range v {
					e.{{ .GoName }} = append(e.{{ .GoName }}, int64(i))
				}
			case []float32:
				e.{{ .GoName }} = make([]int64, 0, len(v))
				for _, i := range v {
					e.{{ .GoName }} = append(e.{{ .GoName }}, int64(i))
				}
			case []float64:
				e.{{ .GoName }} = make([]int64, 0, len(v))
				for _, i := range v {
					e.{{ .GoName }} = append(e.{{ .GoName }}, int64(i))
				}
			case []interface{}:
				e.{{ .GoName }} = make([]int64, 0, len(v))
				for _, i := range v {
					switch vl := i.(type) {
					case int64:
						e.{{ .GoName }} = append(e.{{ .GoName }}, vl)
					case int32:
						e.{{ .GoName }} = append(e.{{ .GoName }}, int64(vl))
					case float32:
						e.{{ .GoName }} = append(e.{{ .GoName }}, int64(vl))
					case float64:
						e.{{ .GoName }} = append(e.{{ .GoName }}, int64(vl))
					default:
            			errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} element parsing %#v: invalid syntax.", i))
					}
				}
			case nil:
			default:
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if eq .Type "[]string" }}
			switch v := value.(type) {
			case []string:
				e.{{ .GoName }} = v
			case []interface{}:
				e.{{ .GoName }} = make([]string, 0, len(v))
				for _, i := range v {
					switch vl := i.(type) {
					case string:
					e.{{ .GoName }} = append(e.{{ .GoName }}, vl)
					default:
						errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} element parsing %#v: invalid syntax.", i))
					}
				}
			case nil:
			default:
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if eq .Type "[]bool" }}
			switch v := value.(type) {
			case []bool:
				e.{{ .GoName }} = v
			case []interface{}:
				e.{{ .GoName }} = make([]bool, 0, len(v))
				for _, i := range v {
					switch vl := i.(type) {
					case bool:
					e.{{ .GoName }} = append(e.{{ .GoName }}, vl)
					default:
						errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} element parsing %#v: invalid syntax.", i))
					}
				}
			case nil:
			default:
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if eq .Type "time.Time" }}
			var v time.Time
			var err error
			valueStr, ok := value.(string)
			if ok {
				switch {
				case constant.NormalDatetimeRegExp.MatchString(valueStr):
					v, err = time.ParseInLocation("2006/01/02 15:04:05", valueStr, time.Local)
				case constant.HyphenDatetimeRegExp.MatchString(valueStr):
					v, err = time.ParseInLocation("2006-01-02 15:04:05", valueStr, time.Local)
				case valueStr == "":
				default:
					v, err = time.Parse(time.RFC3339, valueStr)
				}
				if err != nil {
					errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", valueStr))
				} else {
					e.{{ .GoName }} = v
				}
			} else {
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if eq .Type "[]byte" }}
			switch val := value.(type) {
				case []byte:
					e.{{ .GoName }} = val
				case string:
					v, err := base64.StdEncoding.DecodeString(val)
					if err != nil {
						errs = append(errs, fmt.Sprintf("{{ .GoName }}: []byte base64 decoding %#v: invalid syntax.", value))
					} else {
						e.{{ .GoName }} = v
					}
				default:
					errs = append(errs, fmt.Sprintf("{{ .GoName }}: []byte parsing %#v: invalid syntax.", value))
			}
			{{- else if and .IsList .IsEnum }}
			switch valueStrs := value.(type) {
			case []interface{}:
				e.{{ .GoName }} = make({{ .Type }}, 0, len(valueStrs))
				for _, t := range valueStrs {
					switch vl := t.(type) {
					case string:
						var x {{ trimSuffix "Slice" .Type }}
						err := x.UnmarshalJSON([]byte(fmt.Sprintf("%#v", vl)))
						if err != nil || (x == 0 && vl != "") {
							errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", vl))
						} else {
							e.{{ .GoName }} = append(e.{{ .GoName }}, x)
						}
					default:
						errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} element parsing %#v: invalid syntax.", vl))
					}
				}
			case nil:
			default:
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else if hasPrefix "enum." .Type }}
			valueStr, ok := value.(string)
			if ok {
				err := e.{{ .GoName }}.UnmarshalJSON([]byte(fmt.Sprintf("%#v", valueStr)))
				if err != nil || (e.{{ .GoName }} == 0 && value != "") {
					errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", valueStr))
				}
			} else {
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- else }}
			var ok bool
			e.{{ .GoName }}, ok = value.({{ .Type }})
			if !ok {
				errs = append(errs, fmt.Sprintf("{{ .GoName }}: {{ .Type }} parsing %#v: invalid syntax.", value))
			}
			{{- end }}
		{{- end }}
		}
	}
	return errs
}

// ShallowCopy CreatedTime, UpdatedTime 以外をShallowCopy
func (e *{{ .GoName }}) ShallowCopy() *{{ .GoName }} {
	return &{{ .GoName }}{
		{{- range .Columns }}
		{{ if and (ne "CreatedTime" .GoName) (ne "UpdatedTime" .GoName) }}{{ .GoName }}: e.{{ .GoName }},
		{{- end }}
		{{- end }}
	}
}

// DeepCopy CreatedTime, UpdatedTime 以外をDeepCopy
func (e *{{ .GoName }}) DeepCopy() *{{ .GoName }} {
	{{- range .Columns }}
	{{- if and (not .PK) (.IsList) }}
	{{ .CamelName }} := make({{ .Type }}, len(e.{{ .GoName }}))
	copy({{ .CamelName }}, e.{{ .GoName }})
	{{- end }}
	{{- end }}
	return &{{ .GoName }}{
		{{- range .Columns }}
		{{ if .IsList }}{{ .GoName }}: {{ .CamelName }},
		{{- else if and (ne "CreatedTime" .GoName) (ne "UpdatedTime" .GoName) -}}{{ .GoName }}: e.{{ .GoName }},
		{{- end }}
		{{- end }}
	}
}

// FullDeepCopy 全フィールドをDeepCopy
func (e *{{ .GoName }}) FullDeepCopy() *{{ .GoName }} {
	{{- range .Columns }}
	{{- if and (not .PK) (.IsList) }}
	{{ .CamelName }} := make({{ .Type }}, len(e.{{ .GoName }}))
	copy({{ .CamelName }}, e.{{ .GoName }})
	{{- end }}
	{{- end }}
	return &{{ .GoName }}{
		{{- range .Columns }}
		{{ if .IsList }}{{ .GoName }}: {{ .CamelName }},
		{{- else -}}{{ .GoName }}: e.{{ .GoName }},
		{{- end }}
		{{- end }}
	}
}

// Reset PK, CreatedTime, UpdatedTime 以外を初期化
func (e *{{ .GoName }}) Reset() {
	{{- range .Columns }}
	{{- if and (not .PK) (ne "CreatedTime" .GoName) (ne "UpdatedTime" .GoName) }}
	{{ if eq "string" .Type }}e.{{ .GoName }} = ""
	{{- else if eq "bool" .Type -}}e.{{ .GoName }} = false
	{{- else if eq "time.Time" .Type -}}e.{{ .GoName }} = time.Time{}
	{{- else if or (.IsList) (eq "[]byte" .Type) -}}e.{{ .GoName }} = nil
	{{- else -}}e.{{ .GoName }} = 0
	{{- end }}
	{{- end }}
	{{- end }}
}

{{ if .HasUserID -}}
func (e *{{ .GoName }}) GetUserID() string {
	return e.UserID
}
{{ end -}}

type {{ .GoName }}Slice []*{{ .GoName }}

func (s {{ .GoName }}Slice) GetPKs() {{ .GoName }}PKs {
	pks := make({{ .GoName }}PKs, 0, len(s))
	for _, e := range s {
		pks = append(pks, e.GetPK())
	}
	return pks
}

func (s {{ .GoName }}Slice) Len() int {
	return len(s)
}

func (s {{ .GoName }}Slice) EachRecord(iterator func(Entity) bool) {
	for _, e := range s {
		if !iterator(e) {
			break
		}
	}
}

{{ if and .HasUserID (ne (len .PKColumns) 1) -}}
{{- $pkColumnsOtherThanUserID := slice .PKColumns 1 (len .PKColumns) }}
func (s {{ .GoName }}Slice) CreateUserMap() {{ .GoName }}UserMap {
	m := make({{ .GoName }}UserMap, len(s))
	for _, e := range s {
	{{- range $i, $_ := slice $pkColumnsOtherThanUserID 0 (sub (len $pkColumnsOtherThanUserID) 1) }}
	{{- $cols := slice $pkColumnsOtherThanUserID 0 (add1 $i) }}
	if _, ok := m{{ range $cols }}[e.{{ .GoName }}]{{ end }}; !ok {
		m{{ range slice $cols }}[e.{{ .GoName }}]{{ end }} = make({{ range slice $pkColumnsOtherThanUserID (add1 $i) (len $pkColumnsOtherThanUserID)}}map[{{ .Type }}]{{ end }}*{{ $goName }})
	}
	{{- end }}
	m{{ range $pkColumnsOtherThanUserID }}[e.{{ .GoName }}]{{ end }} = e
	}
	return m
}

type {{ .GoName }}UserMap {{ range $pkColumnsOtherThanUserID }}map[{{ .Type }}]{{ end }}*{{ .GoName }}
{{ end }}

type {{ .GoName }}PK struct {
	{{ range .PKColumns -}}
		{{ .GoName }} {{ .Type }}
	{{ end -}}
}

func (e *{{ .GoName }}PK) ToKeyValue() map[string]interface{} {
	return map[string]interface{}{
	{{- range .PKColumns }}
		"{{ .GoName }}": e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}PK) Generate() []interface{} {
	return []interface{}{
	{{- range .PKColumns }}
		e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}PK) Key() string {
	return strings.Join([]string{
	{{- range .PKColumns }}
		{{ if eq "string" .Type }}e.{{ .GoName }},
		{{- else if eq "int64" .Type -}}strconv.FormatInt(e.{{ .GoName }}, 10),
		{{- else if .IsEnum -}}strconv.FormatInt(e.{{ .GoName }}.Int64(), 10),
		{{- else -}}fmt.Sprint(e.{{ .GoName }}),
		{{ end }}
	{{- end }}
	}, "$")
}

func (e *{{ .GoName }}PK) String() string {
	str := strings.Join([]string{
	{{- range .PKColumns }}
		"{{ .GoName }}: " +
		{{- if eq "string" .Type }}`"` + e.{{ .GoName }} + `"`,
		{{- else if eq "int64" .Type -}}strconv.FormatInt(e.{{ .GoName }}, 10),
		{{- else if .IsEnum -}}`"` + e.{{ .GoName }}.String() + `"`,
		{{- else -}}fmt.Sprint(e.{{ .GoName }}),
		{{ end }}
	{{- end }}
	}, ", ")
	return "{" + str + "}"
}

func (e *{{ .GoName }}PK) ToEntity() *{{ .GoName }} {
	return &{{ .GoName }}{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

{{ if .HasUserID -}}
func (e *{{ .GoName }}PK) GetUserID() string {
	return e.UserID
}
{{- end }}

type {{ .GoName }}PKs []*{{ .GoName }}PK

func (pks {{ .GoName }}PKs) String() string {
	switch len(pks) {
	case 0:
		return "[]"
	case 1:
		return "[" + pks[0].String() + "]"
	}

	n := len(pks) + 1 // 前後の[]とセパレータの空文字分の長さを初期値に
	strs := make([]string, 0, len(pks))
	for _, pk := range pks {
		str := pk.String()
		n += len(str)
		strs = append(strs, str)
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString("[")
	b.WriteString(strs[0])
	for _, str := range strs[1:] {
		b.WriteString(" ")
		b.WriteString(str)
	}
	b.WriteString("]")
	return b.String()
}
{{- range .Indexes }}

// {{ .GoName }} {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .Type }}
{{- end }}
}

func (e *{{ .GoName }}) ToEntity() *{{ $goName }} {
	return &{{ $goName }}{
	{{- range .Columns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

type {{ .GoName }}Slice = []*{{ .GoName }}
{{- end }}

var {{ .GoName }}Columns = column.Columns{
	{{- range .Columns }}
	{{- if not .IsOnlyServer }}
	{
		Name:         "{{ .GoName }}",
		Type:         "{{ .Type }}",
		DatabaseType: "{{ .DatabaseType }}",
		PK:           {{ .PK }},
		Nullable:     false,
		Comment:      "{{ .Comment }}",
	},
	{{- end }}
	{{- end }}
}

var {{ .GoName }}ColumnNames = struct{
{{- range .Columns }}
{{- if not .IsOnlyServer }}
	{{ .GoName }} string
{{- end }}
{{- end }}
}{
{{- range .Columns }}
{{- if not .IsOnlyServer }}
	{{ .GoName }}: "{{ .GoName }}",
{{- end }}
{{- end }}
}

{{ $name := .GoName -}}
var {{ .GoName }}ColumnNameSlice = []string{
{{- range .Columns }}
{{- if not .IsOnlyServer }}
	{{ $name }}ColumnNames.{{ .GoName }},
{{- end }}
{{- end }}
}

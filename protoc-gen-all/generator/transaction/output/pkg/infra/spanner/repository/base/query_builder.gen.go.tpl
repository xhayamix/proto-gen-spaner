{{ template "autogen_comment" }}
{{ $pkgName := .PkgName -}}
{{ $name := .GoName -}}
{{ $camelName := .CamelName -}}
{{ $columns := .Columns -}}
package base

import (
	"strconv"
	"strings"
	"time"

	"github.com/scylladb/go-set/f64set"
	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/i64set"
	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
)


type {{ .GoName }}QueryBuilder interface {
{{- range .Types }}
	Select{{ .GoName }}From{{ $name }}() {{ $name }}QueryBuilderFirstClause
{{- end }}
	SelectCountFrom{{ .GoName }}() {{ .GoName }}QueryBuilderFirstClause
}

type {{ .GoName }}QueryBuilderFinisher interface {
	OrderBy(orderPairs OrderPairs) {{ .GoName }}QueryBuilderFinisher
	Limit(limit int32) {{ .GoName }}QueryBuilderFinisher
	Offset(offset int32) {{ .GoName }}QueryBuilderFinisher
	GetQuery() (string, map[string]any)
	GetQueryConditions() []*{{ .GoName }}QueryCondition
}

type {{ .GoName }}QueryBuilderFirstClause interface {
	{{ .GoName }}QueryBuilderFinisher
	Where() {{ .GoName }}QueryBuilderPredicate
}

type {{ .GoName }}QueryBuilderSecondClause interface {
	{{ .GoName }}QueryBuilderFinisher
	And() {{ .GoName }}QueryBuilderPredicate
}

type {{ .GoName }}QueryBuilderPredicate interface {
	{{ range .Columns -}}
	{{ .GoName }}Eq(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Ne(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Gt(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Gte(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Lt(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Lte(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ if and (ne .SetType "" ) (not .IsList) -}}
	{{ if .IsEnum -}}
	{{ .GoName }}In(params {{ .Type }}Slice) {{ $name }}QueryBuilderSecondClause
	{{ else -}}
	{{ .GoName }}In(params []{{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ .GoName }}Nin(params []{{ .Type }}) {{ $name }}QueryBuilderSecondClause
	{{ end -}}
	{{ end -}}
	{{ end -}}
}

type {{ .GoName }}QueryCondition struct {
	column   string
	operator ConditionOperator
	value    any
}

type {{ .CamelName }}QueryBuilder struct {
	builder         *strings.Builder
	params          map[string]any
	paramIndex      int
	queryConditions []*{{ .GoName }}QueryCondition
}

func New{{ .GoName }}QueryBuilder() {{ .GoName }}QueryBuilder {
	return &{{ .CamelName }}QueryBuilder{
		builder: 		 &strings.Builder{},
		params: 		 make(map[string]any),
		paramIndex: 	 0,
		queryConditions: make( []*{{ .GoName }}QueryCondition, 0),
	}
}

func (qb *{{ $camelName }}QueryBuilder) addParam(condition string, param interface{}) {
	qb.paramIndex++
	paramKey := ParamBaseKey + strconv.Itoa(qb.paramIndex)
	qb.params[paramKey] = param
	qb.builder.WriteString(condition + "@" + paramKey)
}

{{ range .Types -}}
{{ $isIdx := hasPrefix "Idx" .Key }}
func (qb *{{ $camelName }}QueryBuilder) Select{{ .GoName }}From{{ $name }}() {{ $name }}QueryBuilderFirstClause {
	cols := "{{ range $i, $col := .Columns }}{{ if $i }}, {{ end }}`{{ $col.GoName }}`{{ end }}"
	{{- if $isIdx }}
	targetIdx := "@{FORCE_INDEX={{ .Key }}}"
	{{- end }}
	qb.builder.WriteString("SELECT " + cols + " FROM " + "`" + transaction.{{ $name }}TableName + "`"{{- if $isIdx }} + targetIdx{{- end }})
	return qb
}
{{ end -}}

func (qb *{{ .CamelName }}QueryBuilder) SelectCountFrom{{ .GoName }}() {{ .GoName }}QueryBuilderFirstClause {
	qb.builder.WriteString("SELECT COUNT(*) FROM " + "`" + transaction.{{ .GoName }}TableName+ "`")
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Where() {{ .GoName }}QueryBuilderPredicate {
	qb.builder.WriteString(" WHERE ")
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) And() {{ .GoName }}QueryBuilderPredicate {
	qb.builder.WriteString(" AND ")
	return qb
}

{{ range $columns -}}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Eq(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.queryConditions = append(qb.queryConditions, &{{ $name }}QueryCondition{column: "{{ .GoName }}", operator: ConditionOperatorEq, value: param})
	qb.addParam("`{{ .GoName }}` = ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Ne(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.addParam("`{{ .GoName }}` != ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Gt(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.addParam("`{{ .GoName }}` >", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Gte(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.addParam("`{{ .GoName }}` >= ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Lt(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.addParam("`{{ .GoName }}` < ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Lte(param {{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.addParam("`{{ .GoName }}` <= ", param)
	return qb
}

{{ if and (ne .SetType "" ) (not .IsList) -}}
{{- if .IsEnum }}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}In(params {{ .Type }}Slice) {{ $name }}QueryBuilderSecondClause {
	qb.queryConditions = append(qb.queryConditions, &{{ $name }}QueryCondition{column: "{{ .GoName }}", operator: ConditionOperatorIn, value: params.Set()})
{{- else }}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}In(params []{{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	{{- if eq .Type "time.Time" }}
	v := i64set.New()
	for _, t := range params {
		v.Add(t.UnixNano())
	}
	qb.queryConditions = append(qb.queryConditions, &{{ $name }}QueryCondition{column: "{{ .GoName }}", operator: ConditionOperatorIn, value: v})
	{{- else }}
	qb.queryConditions = append(qb.queryConditions, &{{ $name }}QueryCondition{column: "{{ .GoName }}", operator: ConditionOperatorIn, value: {{ .SetType }}.New(params...)})
	{{- end }}
{{- end }}
	qb.builder.WriteString("`{{ .GoName }}` IN (")
	for i, param := range params {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.addParam("", param)
	}
	qb.builder.WriteString(")")
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Nin(params []{{ .Type }}) {{ $name }}QueryBuilderSecondClause {
	qb.builder.WriteString("`{{ .GoName }}` NOT IN (")
	for i, param := range params {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.addParam("", param)
	}
	qb.builder.WriteString(")")
	return qb
}
{{ end }}
{{ end }}

func (qb *{{ .CamelName }}QueryBuilder) OrderBy(orderPairs OrderPairs) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" ORDER BY ")
	for i, pair := range orderPairs {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.builder.WriteString("`" + pair.Column + "` " + string(pair.OrderType))
	}
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Limit(limit int32) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" LIMIT " + strconv.Itoa(int(limit)))
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Offset(offset int32) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" OFFSET " + strconv.Itoa(int(offset)))
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) GetQuery() (string, map[string]any) {
	return qb.builder.String(), qb.params
}

func (qb *{{ .CamelName }}QueryBuilder) GetQueryConditions() []*{{ .GoName }}QueryCondition {
	return qb.queryConditions
}

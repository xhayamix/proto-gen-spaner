package dto

import (
	"database/sql"
	"fmt"

	"cloud.google.com/go/spanner"
)

type MySQLConfig struct {
	User     string
	Password string
	Addr     string
	DB       string
}

func (c *MySQLConfig) GenerateURL() string {
	return fmt.Sprintf(
		"mysql://%s:%s@tcp(%s)/%s",
		c.User,
		c.Password,
		c.Addr,
		c.DB,
	)
}

type SQLSchemaColumn struct {
	TableCatalog           sql.NullString `json:"TABLE_CATALOG"`
	TableSchema            sql.NullString `json:"TABLE_SCHEMA"`
	TableName              sql.NullString `json:"TABLE_NAME"`
	ColumnName             sql.NullString `json:"COLUMN_NAME"`
	OrdinalPosition        sql.NullInt64  `json:"ORDINAL_POSITION"`
	ColumnDefault          sql.NullString `json:"COLUMN_DEFAULT"`
	IsNullable             sql.NullString `json:"IS_NULLABLE"`
	DataType               sql.NullString `json:"DATA_TYPE"`
	CharacterMaximumLength sql.NullInt64  `json:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   sql.NullInt64  `json:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       sql.NullInt64  `json:"NUMERIC_PRECISION"`
	NumericScale           sql.NullInt64  `json:"NUMERIC_SCALE"`
	DatetimePrecision      sql.NullInt64  `json:"DATETIME_PRECISION"`
	CharacterSet           sql.NullString `json:"CHARACTER_SET_NAME"`
	CollationName          sql.NullString `json:"COLLATION_NAME"`
	ColumnType             sql.NullString `json:"COLUMN_TYPE"`
	ColumnKey              sql.NullString `json:"COLUMN_KEY"`
	Extra                  sql.NullString `json:"EXTRA"`
	Privileges             sql.NullString `json:"PRIVILEGES"`
	ColumnComment          sql.NullString `json:"COLUMN_COMMENT"`
	GenerationExpression   sql.NullString `json:"GENERATION_EXPRESSION"`
	SrsID                  sql.NullInt64  `json:"SRS_ID"`
}

type SpannerConfig struct {
	ProjectID string
	Instance  string
	DB        string
}

func (c *SpannerConfig) GenerateURL() string {
	return fmt.Sprintf(
		"spanner://projects/%s/instances/%s/databases/%s?x-clean-statements=true",
		c.ProjectID,
		c.Instance,
		c.DB,
	)
}

type SpannerSchemaColumn struct {
	TableCatalog    spanner.NullString `json:"TABLE_CATALOG"`
	TableSchema     spanner.NullString `json:"TABLE_SCHEMA"`
	TableName       spanner.NullString `json:"TABLE_NAME"`
	ColumnName      spanner.NullString `json:"COLUMN_NAME"`
	OrdinalPosition spanner.NullInt64  `json:"ORDINAL_POSITION"`
	ColumnDefault   spanner.NullString `json:"COLUMN_DEFAULT"`
	DataType        spanner.NullString `json:"DATA_TYPE"`
	IsNullable      spanner.NullString `json:"IS_NULLABLE"`
	SpannerType     spanner.NullString `json:"SPANNER_TYPE"`
}

type SpannerSchemaTable struct {
	TableCatalog                spanner.NullString `json:"TABLE_CATALOG"`
	TableSchema                 spanner.NullString `json:"TABLE_SCHEMA"`
	TableName                   spanner.NullString `json:"TABLE_NAME"`
	ParentTableName             spanner.NullString `json:"PARENT_TABLE_NAME"`
	OnDeleteAction              spanner.NullString `json:"ON_DELETE_ACTION"`
	TableType                   spanner.NullString `json:"TABLE_TYPE"`
	SpannerState                spanner.NullString `json:"SPANNER_STATE"`
	InterleaveType              spanner.NullString `json:"INTERLEAVE_TYPE"`
	RowDeletionPolicyExpression spanner.NullString `json:"ROW_DELETION_POLICY_EXPRESSION"`
}

package output

import (
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
)

type TemplateInfo struct {
	Data     []byte
	FilePath string
}

type EachTemplateCreator interface {
	Create(message *input.Message) (*TemplateInfo, error)
}

type BulkTemplateCreator interface {
	Create(messages []*input.Message, enumMap map[string]*input.Enum) (*TemplateInfo, error)
}

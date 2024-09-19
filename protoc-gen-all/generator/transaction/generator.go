package transaction

import (
	_ "embed"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	migrations "github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/db/ddl/transaction"
	entity "github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/domain/entity/spanner"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/domain/entity/transaction"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/domain/proto/converter"
	repository "github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/domain/repository/spanner"
	repositoryimpl "github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/infra/spanner/repository"
	repositoryimplbase "github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/infra/spanner/repository/base"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output/pkg/registry"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

type generator struct {
	*core.GeneratorBase
	plugin *protogen.Plugin
}

func NewGenerator(plugin *protogen.Plugin) core.Generator {
	return &generator{
		GeneratorBase: core.NewGeneratorBase(),
		plugin:        plugin,
	}
}

var eachCreators = []output.EachTemplateCreator{
	&entity.Creator{},
	&repository.Creator{},
	&repositoryimpl.Creator{},
	&repositoryimplbase.QueryBuilderCreator{},
}

var bulkCreators = []output.BulkTemplateCreator{
	&migrations.Creator{},
	&registry.Creator{},
	&converter.EntityConverterCreator{},
	&converter.PKConverterCreator{},
	&converter.UtilCreator{},
	&transaction.Creator{},
}

func (g *generator) Build() ([]core.GenFile, error) {
	messages := make([]*input.Message, 0)
	enumMap := make(input.EnumMap)

	for _, file := range g.plugin.Files {
		if !file.Generate {
			continue
		}

		if file.Proto.GetPackage() == "server.transaction" {
			message, err := input.ConvertMessageFromProto(file)
			if err != nil {
				return nil, perrors.Stack(err)
			}
			messages = append(messages, message)
			continue
		}

		if file.Proto.GetPackage() == "server.enums" {
			enumMap = enumMap.Merge(input.ConvertEnumsFromProto(file))
			continue
		}
	}

	// 入力ファイルの順番に左右されないようソートする
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].SnakeName < messages[j].SnakeName
	})

	genFiles := make([]core.GenFile, 0)
	for _, creator := range eachCreators {
		for _, message := range messages {
			info, err := creator.Create(message)
			if err != nil {
				return nil, perrors.Stack(err)
			}
			if info == nil {
				continue
			}

			genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
		}
	}
	for _, creator := range bulkCreators {
		info, err := creator.Create(messages, enumMap)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		if info == nil {
			continue
		}

		genFiles = append(genFiles, core.NewGenFile(info.FilePath, info.Data))
	}

	return genFiles, nil
}

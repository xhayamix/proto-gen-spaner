package main

import (
	"strings"
	"time"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/api"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/plogging"
)

type flagKind string

const (
	flagKindGenTransaction flagKind = "gen_transaction"
	flagKindGenEnum        flagKind = "gen_enum"
	flagKindGenApi         flagKind = "gen_api"
)

func main() {
	locationName := "Asia/Tokyo"
	location, err := time.LoadLocation(locationName)
	if err != nil {
		location = time.FixedZone(locationName, 9*60*60)
	}
	time.Local = location

	startTime := time.Now()
	logger := plogging.GetLogger()
	logger.Infof("proto-gen-spanner start\n")

	generatorBuilder := core.NewGeneratorBuilder()

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		generatorMap, writePb := createGeneratorMap(plugin)
		kinds := make([]string, 0, len(generatorMap))
		for kind, generator := range generatorMap {
			kinds = append(kinds, string(kind))
			generatorBuilder.AppendGenerator(generator)
		}
		logger.Infof("flag %s\n", strings.Join(kinds, ","))

		generatedFilenamePrefixList := make([]string, 0, len(plugin.Files))
		if writePb {
			for _, file := range plugin.Files {
				if !file.Generate {
					continue
				}
				generatedFilenamePrefixList = append(generatedFilenamePrefixList, file.GeneratedFilenamePrefix)
			}
		}
		if err := generatorBuilder.Generate(generatedFilenamePrefixList); err != nil {
			return perrors.Stack(err)
		}
		return nil
	})

	endTime := time.Now()
	logger.Infof("proto-gen-spanner end, elapsed: %s\n", endTime.Sub(startTime).String())
}

func createGeneratorMap(plugin *protogen.Plugin) (map[flagKind]core.Generator, bool) {
	generatorMap := make(map[flagKind]core.Generator)
	var writePb bool

	for _, param := range strings.Split(plugin.Request.GetParameter(), ",") {
		s := strings.Split(param, "=")

		switch flagKind(s[0]) {
		case flagKindGenTransaction:
			generatorMap[flagKindGenTransaction] = transaction.NewGenerator(plugin)
		case flagKindGenEnum:
			generatorMap[flagKindGenEnum] = enum.NewGenerator(plugin)
		case flagKindGenApi:
			generatorMap[flagKindGenApi] = api.NewGenerator(plugin)
		default:
			continue
		}
	}

	return generatorMap, writePb
}

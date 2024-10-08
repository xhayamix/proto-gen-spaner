{{ template "autogen_comment" }}
syntax = "proto3";

package client.master;
option go_package = "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/master";

import "client/master/common/options.proto";

// {{ .Comment }}
message {{ .PascalPrefix }}Setting {
  option (master.common.table) = true;

  // 設定ID（1レコードのみ"1"で固定）
  string id = 1 [(master.common.pk) = true];

  {{ range $i, $e := .Elements -}}
  // {{ $e.Comment }}
  {{ if .IsList }}repeated {{ end }}{{ $e.SettingType }} {{ $e.CamelName }} = {{ $e.Value }};
  {{ end -}}
}

message {{ .PascalPrefix }}SettingList {
  repeated {{ .PascalPrefix }}Setting list = 1;
}

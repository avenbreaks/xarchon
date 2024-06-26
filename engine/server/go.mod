module github.com/avenbreaks/xarchon/engine/server

go 1.18

replace (
	github.com/avenbreaks/xarchon/api/golang => ../../api/golang
	github.com/avenbreaks/xarchon/container-engine-lib => ../../container-engine-lib
	github.com/avenbreaks/xarchon/contexts-config-store => ../../contexts-config-store
	github.com/avenbreaks/xarchon/grpc-file-transfer/golang => ../../grpc-file-transfer/golang
	github.com/avenbreaks/xarchon/core/launcher => ../../core/launcher
	github.com/avenbreaks/xarchon/engine/launcher => ../launcher
	github.com/avenbreaks/xarchon/xarchon_version => ../../xarchon_version
	github.com/avenbreaks/xarchon/name_generator => ../../name_generator
)

require (
	github.com/avenbreaks/xarchon/api/golang v0.0.0
	github.com/avenbreaks/xarchon/container-engine-lib v0.0.0 // local dependency
	github.com/avenbreaks/xarchon/core/launcher v0.0.0 // local dependency
	github.com/avenbreaks/xarchon/engine/launcher v0.0.0
	github.com/avenbreaks/xarchon/name_generator v0.0.0 // local dependency
	github.com/kurtosis-tech/minimal-grpc-server/golang v0.0.0-20211201000847-a204edc5a0b3
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.29.1
)

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.0+incompatible // indirect
	github.com/docker/docker v20.10.24+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gammazero/deque v0.1.0 // indirect
	github.com/gammazero/workerpool v1.1.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20210402141018-6c239bbf2bb1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/gorilla/websocket v1.4.2

require (
	github.com/avenbreaks/xarchon/xarchon_version v0.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

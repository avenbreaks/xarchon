module github.com/avenbreaks/xarchon/core/files_artifacts_expander

go 1.18

replace (
	github.com/avenbreaks/xarchon/api/golang => ../../api/golang
	github.com/avenbreaks/xarchon/contexts-config-store => ../../contexts-config-store
	github.com/avenbreaks/xarchon/grpc-file-transfer/golang => ../../grpc-file-transfer/golang
)

require (
	github.com/avenbreaks/xarchon/api/golang v0.0.0 // Local dependency
	github.com/avenbreaks/xarchon/grpc-file-transfer/golang v0.0.0 // Local dependency
	github.com/gammazero/workerpool v1.1.2
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.53.0
)

require (
	github.com/gammazero/deque v0.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

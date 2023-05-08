module github.com/avenbreaks/xarchon/core/launcher

go 1.18

replace (
	github.com/avenbreaks/xarchon/container-engine-lib => ../../container-engine-lib
	github.com/avenbreaks/xarchon/xarchon_version => ../../xarchon_version
)

require (
	github.com/avenbreaks/xarchon/container-engine-lib v0.0.0 // Local dependency
	github.com/avenbreaks/xarchon/xarchon_version v0.0.0 // Local dependency generated during build
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

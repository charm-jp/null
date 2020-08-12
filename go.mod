module github.com/charm-jp/null

//replace gitlab.dev.charm.internal/grpc/grpc-shared-go => ../../../gitlab.dev.charm.internal/grpc/grpc-shared-go

go 1.12

require (
	gitlab.dev.charm.internal/grpc/grpc-shared-go v0.0.0-20200811015345-d2c5fb6b7374
	google.golang.org/protobuf v1.25.0 // indirect
)

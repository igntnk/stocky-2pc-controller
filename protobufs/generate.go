package protobufs

//go:generate protoc --go_out=./iims_pb/ --go_opt=paths=source_relative --go-grpc_out=./iims_pb/ --go-grpc_opt=paths=source_relative iims.proto
//go:generate protoc --go_out=./sms_pb/ --go_opt=paths=source_relative --go-grpc_out=./sms_pb/ --go-grpc_opt=paths=source_relative sms.proto
//go:generate protoc --go_out=./scs_pb/ --go_opt=paths=source_relative --go-grpc_out=./scs_pb/ --go-grpc_opt=paths=source_relative scs.proto
//go:generate protoc --go_out=./oms_pb/ --go_opt=paths=source_relative --go-grpc_out=./oms_pb/ --go-grpc_opt=paths=source_relative oms.proto

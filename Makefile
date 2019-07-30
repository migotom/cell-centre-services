.PHONY: proto

proto:
	protoc -I pkg/pb/ -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		pkg/pb/role.proto \
		pkg/pb/employee.proto \
		pkg/pb/auth.proto \
		pkg/pb/event.proto \
		--go_out=plugins=grpc:pkg/pb
		
	protoc -I pkg/pb/ -I$$GOPATH/src -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		pkg/pb/employee.proto \
		--grpc-gateway_out=logtostderr=true,grpc_api_configuration=pkg/pb/employee_service.yaml:pkg/pb
	protoc -I pkg/pb/ -I$$GOPATH/src -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		pkg/pb/auth.proto \
		--grpc-gateway_out=logtostderr=true,grpc_api_configuration=pkg/pb/auth_service.yaml:pkg/pb
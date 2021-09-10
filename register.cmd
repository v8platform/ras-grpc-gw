
grpcurl -vv -protoset ./protos/server.bin -plaintext -d "{ \"user\": \"admin\", \"password\": \"pwd\"}" localhost:3002 access.service.AuthService/SingIn
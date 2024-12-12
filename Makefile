
.PHONY: build
build: build-buf build-migrate build-sqlc-pgx

# builds the proto files
.PHONY: build-buf
build-buf:
	buf dep update
	buf generate -v

.PHONY: build-migrate
build-migrate:
	migrate create -ext sql -dir migrations -seq init

.PHONY: build-sqlc
build-sqlc:
	sqlc generate

.PHONY: build-sqlc-pgx
build-sqlc-pgx:
	sqlc generate -f sqlc.pgx.yaml

# stop the docker-compose
.PHONY: down
down:
	docker-compose down

# start docker-compose
.PHONY: up
up:
	. .env
	docker-compose up --build

# deletes all docker contents
.PHONY: clean
clean:
	docker system prune -a -f

# cert commands for tls
.PHONY: gen-certs
gen-certs: gen-ca-certs gen-server-certs gen-client-certs

# The CA acts as a trusted third party that confirms the identities of all parties involved.
.PHONY: gen-ca-certs
gen-ca-certs:
	mkdir -p certs
	openssl genrsa -out certs/ca.key 4096
	openssl req -x509 -new -nodes -key certs/ca.key -sha256 -days 1024 -out certs/ca.crt -subj "/CN=gRPC CA"

.PHONY: gen-server-certs
gen-server-certs:
	openssl genrsa -out certs/server.key 4096
	openssl req -new -key certs/server.key -out certs/server.csr -subj "/CN=gRPC TLS Server Example"
	echo "subjectAltName=IP:0.0.0.0,IP:127.0.0.1,DNS:localhost" > certs/server_extfile.cnf
	openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt -days 500 -sha256 -extfile certs/server_extfile.cnf

.PHONY: gen-client-certs
gen-client-certs:
	openssl genrsa -out certs/client.key 4096
	openssl req -new -key certs/client.key -out certs/client.csr -subj "/CN=gRPC TLS Client Example"
	echo "subjectAltName=IP:0.0.0.0,IP:127.0.0.1,DNS:localhost" > certs/client_extfile.cnf
	openssl x509 -req -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt -days 500 -sha256 -extfile certs/client_extfile.cnf

# Below tests configured for a server side tls
# for insecure use -plaintext flag
# If testing with server running over tls replace -plaintext with -cacert certs/ca.crt
# If testing with server running over mtls replace -plaintext with -cacert certs/ca.crt -cert certs/client.crt -key certs/client.key
.PHONY: test-get-grpc
test-get-grpc:
	grpcurl -cacert certs/ca.crt \
	  localhost:9090 user.v1.UserService/GetUsers

.PHONY: test-post-grpc
test-post-grpc:
	grpcurl -cacert certs/ca.crt \
	  -d '{"user": {"username":"new user", "email": "new email", "first_name": "new", "last_name": "user"}}' \
	  localhost:9090 user.v1.UserService/CreateUser

.PHONY: test-update-grpc
test-update-grpc:
	grpcurl -cacert certs/ca.crt \
	  -d '{"user_id": 4, "username":"updated user", "email": "updated email", "first_name": "updated", "last_name": "user"}' \
	  localhost:9090 user.v1.UserService/UpdateUser

.PHONY: test-delete-grpc
test-delete-grpc:
	grpcurl -cacert certs/ca.crt \
	  -d '{"user_id": 4}' \
	  localhost:9090 user.v1.UserService/DeleteUser

.PHONY: test-get-rest
test-get-rest:
	curl -H "Content-Type: application/json" \
	  --cacert certs/ca.crt \
	  https://localhost:8080/api/v1/users

.PHONY: test-stream-rest
test-stream-rest:
	curl -H "Content-Type: application/json" \
	  --cacert certs/ca.crt \
	  https://localhost:8080/api/v1/userstream

.PHONY: test-post-rest
test-post-rest:
	curl -X POST -H "Content-Type: application/json" \
	  --cacert certs/ca.crt \
	  -d '{"user": {"username":"new user", "email": "new email", "first_name": "new", "last_name": "user"}}' \
          https://localhost:8080/api/v1/users

.PHONY: test-update-rest
test-update-rest:
	curl -X PUT -H "Content-Type: application/json" \
	  --cacert certs/ca.crt \
	  -d '{"user_id": 5, "username":"updated user", "email": "updated email", "first_name": "updated", "last_name": "user"}' \
          https://localhost:8080/api/v1/users

.PHONY: test-delete-rest
test-delete-rest:
	curl -X DELETE -H "Content-Type: application/json" \
	  --cacert certs/ca.crt \
          https://localhost:8080/api/v1/users/5



# builds the proto files
.PHONY: build
build:
	buf dep update
	buf generate -v

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
	# mkdir -p certs
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
	grpcurl -cacert certs/ca.crt localhost:9090 user.v1.UserService/GetUsers

.PHONY: test-post-grpc
test-post-grpc:
	grpcurl -cacert certs/ca.crt \
	-d '{"user": {"Username":"new user", "Email": "new email", "FirstName": "new", "LastName": "user"}}' \
	localhost:9090 user.v1.UserService/CreateUser

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
          -d '{"Username":"new user", "Email": "new email", "FirstName": "new", "LastName": "user"}' \
          https://localhost:8080/api/v1/users
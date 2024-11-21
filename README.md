# go-app-tmpl
A template for generic golang applications.

## Tooling

This section will cover all the moving parts of this demo. Each closes one gap for the fully declarative extensible build.

### [sqlc]("https://sqlc.dev/")

`sqlc` is used to build out your database connection and query layers using just regular sql queries. 

### [buf]("https://github.com/bufbuild/buf")

`buf` is the new industry standard when it comes to protobuf generation. `protoc` is still used but `buf` is more declarative and extensible.

### [golang-migrate]("https://github.com/golang-migrate/migrate")

golang-migrate `migrate` tool is used for database migrations. Its somewhat overkill for this but nice to show that it can be used along with `sqlc` and is easier to manage sql changes than traditional sql scripts.

## Initialization

1. Install the following tools.

- sqlc cli

- buf cli

- golang-migrate cli

- docker and docker-compose

- make (optional)

- grpcurl (optional for testing grpc server)

2. Setup `.env` file - nothing special just to setup db creds.

```bash
export DB_HOST=localhost
export DB_PASSWORD=your_password
export DB_USER=your_username
export DB_NAME=golang_app
```

**Note: the following are not needed if cloning this repo but ill cover for completeness.**

3. Initialize the go module.

```bash
go mod init github.com/alexander-bergeron/go-app-tmpl
```

4. Initialize Docker files. 

_see file contents for further explanation_

- `Dockerfile.server` builds the golang server executable.

- `compose.yml` is the docker-compose file for orchestration.

5. Initialize migrations

```bash
migrate create -ext sql -dir migrations -seq init
```

This creates your migrations directory and creates two empty files, init.up and init.down.

```sql
-- init.up
-- Create the users table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(100),
    last_name VARCHAR(100)
);

INSERT INTO users (username, email, first_name, last_name) VALUES
('Cow31337Killer', 'ckilla@hotmail.com', 'Cow', 'Killer'),
('Durial321', 'backslash@yahoo.com', 'Durial', '321'),
('BigRedJapan', 'brj@gmail.com', 'BigRed', 'Japan');
```

```sql
-- init.down
DROP TABLE IF EXISTS users;
```

6. Initialize `sqlc.pgx.yaml`.

_see file contents for further explanation_

7. Initialize `buf` files.

_see file contents for further explanation_

## Build Steps

1. `make build`

- This will generate the grpc protobuf code, the migrate directory with init files, and the sqlc.pgx.yaml code.

- After build completes, manually update swagger.json with host and schemes.

```json
  "info": {
    "title": "proto/user/v1/user.proto",
    "version": "version not set"
  },
  "host": "localhost:8080", # add this line
  "schemes": ["https"], # add this line
  "tags": [
    {
      "name": "UserService"
    }
  ],
```

2. `go mod tidy`

- Updates all your deps.

3. `make gen-certs`

- Creates cert files for tls.

4. start docker daemon.

5. `make up`

- Runs docker compose.

## Testing

1. `make test-get-grpc` - test grpc client call

2. `make test-get-rest` - test grpc-gateway

3. `make test-post-grpc` - posts a new test user to the API

4. Navigate to `localhost:3000` for a webserver displaying all users with form input for a new user - built with htmx

5. Navigate to `localhost:8081` for a swagger page - **Note: you will need to add the generated `ca.crt` to your browser for swagger to work.**

## TLS

Serving over TLS takes various steps.

1. creating the certs - see Makefile

2. Baking ca cert into swagger ui image

3. updating browser to recognize your custom ca - import the ca.crt

## TODO

1. Determine how to set host automatically for swagger page rather than manual change after generation.

2. Setup CORS so its not hard coded.

3. setup static css and js

4. test to see if air works in container - [reference](https://afarid.medium.com/golang-hot-reloading-using-docker-and-air-b6da91293cd9)

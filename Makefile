GOFMT_FILES?=$$(find . -name '*.go')

all: fmt build install

fmt:
	gofmt -w $(GOFMT_FILES)

build: fmt
	go build -ldflags "-X github.com/swiftcarrot/queryx/cmd/queryx/action.Version=`git rev-parse HEAD`" -o bin/queryx cmd/queryx/main.go

install: build
	install bin/queryx /usr/local/bin

clean:
	rm bin/queryx

test-postgresql: install
	rm -rf internal/integration/db
	cd internal/integration && QUERYX_ENV=test queryx db:drop --schema postgresql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:create --schema postgresql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema postgresql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema postgresql.hcl
	cd internal/integration && QUERYX_ENV=test queryx generate --schema postgresql.hcl
	cd internal/integration && go test ./...
	# cd internal/integration && QUERYX_ENV=test queryx db:drop --schema postgresql.hcl

test-mysql: install
	rm -rf internal/integration/db
	cd internal/integration && QUERYX_ENV=test queryx db:drop --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:create --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx generate --schema mysql.hcl
	cd internal/integration && go test ./...

test-sqlite: install
	rm -rf internal/integration/db
	cd internal/integration && QUERYX_ENV=test queryx db:drop --schema sqlite.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:create --schema sqlite.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema sqlite.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema sqlite.hcl
	cd internal/integration && QUERYX_ENV=test queryx generate --schema sqlite.hcl
	cd internal/integration && go test ./...

test: test-postgresql test-sqlite test-mysql

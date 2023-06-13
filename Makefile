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
	cd internal/integration && QUERYX_ENV=test queryx generate --schema postgresql.hcl
	cd internal/integration && go test -v ./...
	# cd internal/integration && QUERYX_ENV=test queryx db:drop --schema postgresql.hcl

test-mysql: install mysql-drop
	rm -rf internal/integration/db/migrations/* && rm -rf internal/integration/db/builder && rm -rf internal/integration/db/*.go
	cd internal/integration && QUERYX_ENV=test queryx generate --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:create --schema mysql.hcl
	cd internal/integration && QUERYX_ENV=test queryx db:migrate --schema mysql.hcl
	cd internal/integration && go test -v ./...

test-sqlite: install
	rm -rf internal/integration/db/migrations/* && rm -rf internal/integration/db/builder && rm -rf internal/integration/db/*.go &&  rm -rf internal/integration/*.db
	cd internal/integration && rm -f test.sqlite3 && touch test.sqlite3
	cd internal/integration && QUERYX_ENV=test DATABASE_URL=test.sqlite3 queryx db:migrate --schema sqlite.hcl
	cd internal/integration && QUERYX_ENV=test DATABASE_URL=test.sqlite3 queryx generate --schema sqlite.hcl
	cd internal/integration && go test -v ./...

sqlite: install
	cd internal/integration && rm -rf db/*.go db/builder && queryx generate --schema sqlite.hcl && rm -f /opt/data.db && touch /opt/data.db
	cd internal/integration && QUERYX_ENV=test DATABASE_URL=/opt/data.db  go test -v ./...

test: test-pg test-mysql test-sqlite

pg-drop:
	cd internal/integration && QUERYX_ENV=test queryx db:drop --schema postgresql.hcl

mysql-drop:
	cd internal/integration && QUERYX_ENV=test queryx db:drop --schema mysql.hcl

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
	cd internal/integration && queryx db:drop --schema postgresql.hcl
	cd internal/integration && queryx db:create --schema postgresql.hcl
	cd internal/integration && queryx db:migrate --schema postgresql.hcl
	cd internal/integration && queryx db:migrate --schema postgresql.hcl
	cd internal/integration && queryx generate --schema postgresql.hcl
	cd internal/integration && yarn tsc
	cd internal/integration && yarn test
	# cd internal/integration && go test ./...
	# cd internal/integration && queryx db:drop --schema postgresql.hcl

test-mysql: install
	rm -rf internal/integration/db
	cd internal/integration && queryx db:drop --schema mysql.hcl
	cd internal/integration && queryx db:create --schema mysql.hcl
	cd internal/integration && queryx db:migrate --schema mysql.hcl
	cd internal/integration && queryx db:migrate --schema mysql.hcl
	cd internal/integration && queryx generate --schema mysql.hcl
	cd internal/integration && yarn tsc
	cd internal/integration && yarn test
	# cd internal/integration && go test ./...

test-sqlite: install
	rm -rf internal/integration/db
	cd internal/integration && queryx db:drop --schema sqlite.hcl
	cd internal/integration && queryx db:create --schema sqlite.hcl
	cd internal/integration && queryx db:migrate --schema sqlite.hcl
	cd internal/integration && queryx db:migrate --schema sqlite.hcl
	cd internal/integration && queryx generate --schema sqlite.hcl
	cd internal/integration && yarn tsc
	cd internal/integration && yarn test
	# cd internal/integration && go test ./...

test: test-postgresql test-sqlite test-mysql

test-migrate: install
	rm -rf internal/migrate/db internal/migrate/test.sqlite3
	cd internal/migrate && queryx db:migrate --schema sqlite1.hcl
	sleep 1
	cd internal/migrate && queryx db:migrate --schema sqlite2.hcl
	cd internal/migrate && sqlite3 test.sqlite3 "insert into users(name, email) values('test', 'test@example.com')"


benchmarks-golang-postgresql:
	cd internal/bench/go-queryx && rm -rf db
	cd internal/bench/go-queryx && queryx db:drop --schema postgresql.hcl
	cd internal/bench/go-queryx && queryx db:create --schema postgresql.hcl
	cd internal/bench/go-queryx && queryx db:migrate --schema postgresql.hcl
	cd internal/bench/go-queryx && queryx g --schema postgresql.hcl
	cd internal/bench && go build -o bin/queryxorm main.go
	cd internal/bench && install bin/queryxorm /usr/local/bin
	queryxorm -orm=go-queryx

benchmarks-golang-mysql:
	cd internal/bench/go-queryx && rm -rf db
	cd internal/bench/go-queryx && queryx db:drop --schema mysql.hcl
	cd internal/bench/go-queryx && queryx db:create --schema mysql.hcl
	cd internal/bench/go-queryx && queryx db:migrate --schema mysql.hcl
	cd internal/bench/go-queryx && queryx g --schema mysql.hcl
	cd internal/bench && go build -o bin/queryxorm main.go
	cd internal/bench && install bin/queryxorm /usr/local/bin
	queryxorm -orm=go-queryx

benchmarks-golang-sqlite:
	cd internal/bench/go-queryx && rm -rf db
	cd internal/bench/go-queryx && queryx db:drop --schema sqlite.hcl
	cd internal/bench/go-queryx && queryx db:create --schema sqlite.hcl
	cd internal/bench/go-queryx && queryx db:migrate --schema sqlite.hcl
	cd internal/bench/go-queryx && queryx g --schema sqlite.hcl
	cd internal/bench && go build -o bin/queryxorm main.go
	cd internal/bench && install bin/queryxorm /usr/local/bin
	queryxorm -orm=go-queryx

benchmarks-typescript-postgresql:
	cd internal/bench/ts-queryx && rm -rf db
	cd internal/bench/ts-queryx && queryx db:drop --schema postgresql.hcl
	cd internal/bench/ts-queryx && queryx db:create --schema postgresql.hcl
	cd internal/bench/ts-queryx && queryx db:migrate --schema postgresql.hcl
	cd internal/bench/ts-queryx && queryx g --schema postgresql.hcl
	cd internal/bench/ts-queryx && tsc benchmark.test.ts
	cd internal/bench/ts-queryx && node benchmark.test.js


benchmarks-typescript-mysql:
	cd internal/bench/ts-queryx && rm -rf db
	cd internal/bench/ts-queryx && queryx db:drop --schema    mysql.hcl
	cd internal/bench/ts-queryx && queryx db:create --schema  mysql.hcl
	cd internal/bench/ts-queryx && queryx db:migrate --schema mysql.hcl
	cd internal/bench/ts-queryx && queryx g --schema          mysql.hcl
	cd internal/bench/ts-queryx &&  yarn
	cd internal/bench/ts-queryx && tsc benchmark.test.ts
	cd internal/bench/ts-queryx && node benchmark.test.js


benchmarks-typescript-sqlite:
	cd internal/bench/ts-queryx && rm -rf db
	cd internal/bench/ts-queryx && queryx db:drop --schema    sqlite.hcl
	cd internal/bench/ts-queryx && queryx db:create --schema  sqlite.hcl
	cd internal/bench/ts-queryx && queryx db:migrate --schema sqlite.hcl
	cd internal/bench/ts-queryx && queryx g --schema          sqlite.hcl
	cd internal/bench/ts-queryx && tsc benchmark.test.ts
	cd internal/bench/ts-queryx && node benchmark.test.js

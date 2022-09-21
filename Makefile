GO111MODULE=on
GOPRIVATE=source.golabs.io
GOSUMDB=sum.golang.org
GOPROXY=proxy.golang.org
export GO111MODULE
export GOPRIVATE
export GOSUMDB
export GOPROXY

# This will install all package dependencies
setup:
	@echo "installing all required dependencies..."
	go env -w GOPRIVATE=source.golabs.io
	go get -d -v
	go mod vendor
	go mod tidy

# This will clean the build directory
clean:
	@echo "cleaning..."
	rm -rf bin
	go mod tidy
	@echo "cleaning done"

# This will generate an executable in ./out directory.
build: setup clean
	@echo "building..."
	mkdir -p bin
	go build -o bin/url-blaster -v .

# This will run golangci-lint
lint:
	@echo "linting using golang-ci lint"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0 run

# This will run all the tests.
test:
	@echo "running tests..."
	go mod tidy
	go test -coverpkg=./... ./...

# This will run the service 
run: build
	@echo "running..."
	./bin/url-blaster

# This will run CI
ci: build lint

# This will run the entire process from building, running, to cleaning again
run.all:
	make run
	make clean

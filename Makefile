# This will install all package dependencies
setup:
	@echo "installing all required dependencies..."
	go get -d -v
	go mod vendor
	go mod tidy

# This will clean the build directory
clean:
	@echo "cleaning..."
	rm -rf bin
	go mod tidy

# This will generate an executable in ./out directory.
build: setup clean
	@echo "building..."
	mkdir -p bin
	go build -o bin/url-blaster -v .

# This will run golangci-lint
lint:
	@echo "linting using golang-ci lint BUT NOT YET IMPLEMENTED"

# This will run all the tests.
test:
	@echo "running tests..."
	go mod tidy
	go test

# This will run the service 
run:
	@echo "running..."
	go run main.go

# This will run CI
ci: build lint
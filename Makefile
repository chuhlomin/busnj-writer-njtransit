build:
	@cd ./cmd/writer; \
	go build .

build-static:
	@cd ./cmd/writer; \
	CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o writer .

run:
	@cd ./cmd/writer; \
	go run .

vet:
	@go vet ./cmd/...

test:
	@go test ./...

build-docker:
	@docker build --tag busnj-writer-njtransit:latest ./cmd/writer;

run-docker:
	@docker run --name busnj-writer-njtransit \
		--rm \
		--network busnj-network \
		--env BUSDATA_USERNAME=${BUSDATA_USERNAME} \
		--env BUSDATA_PASSWORD=${BUSDATA_PASSWORD} \
		busnj-writer-njtransit:latest;

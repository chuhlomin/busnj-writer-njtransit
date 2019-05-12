build:
	@cd ./cmd/writer; \
	go build .

build-drone:
	@cd ./cmd/writer; \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o writer .

run:
	@cd ./cmd/writer; \
	go run .

vet:
	@go vet ./cmd/...

test:
	@go test ./...

docker-build:
	@docker build --tag busnj-writer-njtransit:latest ./cmd/writer;

docker-run:
	@docker run --name busnj-writer-njtransit \
		--rm \
		--network busnj-network \
		--env BUSDATA_USERNAME=${BUSDATA_USERNAME} \
		--env BUSDATA_PASSWORD=${BUSDATA_PASSWORD} \
		--env REDIS_ADDR=${REDIS_ADDR} \
		busnj-writer-njtransit:latest;

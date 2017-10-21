.PHONY : test format

build:
	docker build -t wuriyanto/go-ddd-grpc .

test:
	go test ./server/model \
					./server/query

format:
	find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" | xargs gofmt -s -d -w

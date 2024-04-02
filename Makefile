PROTOCOL_VERSION?=v0.6.0

clean:
	rm protocol/v1/saturnsync.proto

generate: protocol/v1/saturnsync.proto
	buf generate
	protoc-go-inject-tag -input ./protocol/v1/saturnsync.pb.go -remove_tag_comment
	gofmt -w ./protocol/v1/saturnsync.pb.go

protocol/v1/saturnsync.proto:
	curl -L --silent --fail -o ./protocol/v1/saturnsync.proto https://raw.githubusercontent.com/wndhydrnt/saturn-sync-protocol/$(PROTOCOL_VERSION)/protocol/v1/saturnsync.proto

test_cover:
	go test -covermode=set -coverpkg=./... -coverprofile cover.out -v ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out
